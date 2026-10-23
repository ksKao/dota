package lib

import (
	"errors"
	"io"
	"strings"

	"github.com/dotabuff/manta"
	"github.com/dotabuff/manta/dota"
	"github.com/ksKao/dota/models"
)

type eventType int

const (
	eventTypeCast     eventType = 0
	eventTypeDamage   eventType = 1
	eventTypeModifier eventType = 2
)

func getName(p *manta.Parser, index uint32) string {
	var name string
	var ok bool
	if name, ok = p.LookupStringByIndex("CombatLogNames", int32(index)); !ok {
		return "<unknown>"
	}
	return name
}

func ParseDemFile(r io.Reader) (start float32, hooks []models.Hook, err error) {
	p, err := manta.NewStreamParser(r)
	if err != nil {
		return
	}

	p.OnEntity(func(e *manta.Entity, op manta.EntityOp) error {
		if start != 0 || e.GetClassName() != "CDOTAGamerulesProxy" {
			return nil
		}

		start, _ = e.GetFloat32("m_pGameRules.m_flGameStartTime")
		return nil
	})

	type event struct {
		Type      eventType
		Timestamp float32
		Target    string
	}

	events := []event{}
	p.Callbacks.OnCMsgDOTACombatLogEntry(func(m *dota.CMsgDOTACombatLogEntry) error {
		caster := getName(p, m.GetAttackerName())
		ability := getName(p, m.GetInflictorName())
		target := getName(p, m.GetTargetName())
		timestamp := m.GetTimestamp()
		if caster != "npc_dota_hero_pudge" {
			return nil
		}
		if caster != "npc_dota_hero_pudge" || (ability != "pudge_meat_hook" && ability != "modifier_pudge_meat_hook") {
			return nil
		}

		event := event{
			Target:    target,
			Timestamp: timestamp,
		}
		switch m.GetType() {
		case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_ABILITY:
			event.Type = eventTypeCast
		case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_DAMAGE:
			event.Type = eventTypeDamage
		case dota.DOTA_COMBATLOG_TYPES_DOTA_COMBATLOG_MODIFIER_ADD:
			event.Type = eventTypeModifier
		default:
			return nil
		}

		events = append(events, event)
		return nil
	})

	p.Start()

	if len(events) == 0 {
		return
	}

	i := 0
	for i < len(events) {
		event := events[i]

		if event.Type != eventTypeCast {
			err = errors.New("current event is not cast")
			return
		}

		// if there's no element after, missed
		if i == len(events)-1 {
			hooks = append(hooks, *models.NewHook(event.Timestamp).Missed())
			return
		}

		next := events[i+1]
		switch next.Type {
		case eventTypeCast:
			hooks = append(hooks, *models.NewHook(event.Timestamp).Missed())
			i++
		case eventTypeDamage:
			// check if there's an element after next
			if i == len(events)-2 {
				hooks = append(hooks, *models.NewHook(event.Timestamp).NotMissed(&event.Target).Enemy().Creep())
				return
			}
			afterNext := events[i+2]
			switch afterNext.Type {
			case eventTypeCast:
				hooks = append(hooks, *models.NewHook(event.Timestamp).NotMissed(&event.Target).Enemy().Creep())
				i += 2
			case eventTypeModifier:
				hooks = append(hooks, *models.NewHook(event.Timestamp).NotMissed(&event.Target).Enemy().Hero())
				i += 3
			default:
				err = errors.New("unexpected consecutive damage event")
				return
			}
		case eventTypeModifier:
			if strings.HasPrefix(next.Target, "npc_dota_creep") {
				hooks = append(hooks, *models.NewHook(event.Timestamp).NotMissed(&event.Target).Ally().Creep())
			} else {
				hooks = append(hooks, *models.NewHook(event.Timestamp).NotMissed(&event.Target).Ally().Hero())
			}
			i += 2
		}
	}

	return
}
