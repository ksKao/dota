package models

type Hook struct {
	Timestamp float32 `json:"timestamp"`
	IsMissed  bool    `json:"isMissed"`
	IsEnemy   *bool   `json:"isEnemy,omitempty"`
	IsCreep   *bool   `json:"isCreep,omitempty"`
	Target    *string `json:"target,omitempty"`
}

func NewHook(timestamp float32) *Hook {
	return &Hook{
		Timestamp: timestamp,
	}
}

func (h *Hook) Missed() *Hook {
	h.IsMissed = true
	return h
}

func (h *Hook) NotMissed(target *string) *Hook {
	h.IsMissed = false
	h.Target = target
	return h
}

func (h *Hook) Enemy() *Hook {
	isEnemy := true
	h.IsEnemy = &isEnemy
	return h
}

func (h *Hook) Ally() *Hook {
	isEnemy := false
	h.IsEnemy = &isEnemy
	return h
}

func (h *Hook) Creep() *Hook {
	isCreep := true
	h.IsCreep = &isCreep
	return h
}

func (h *Hook) Hero() *Hook {
	isCreep := false
	h.IsCreep = &isCreep
	return h
}
