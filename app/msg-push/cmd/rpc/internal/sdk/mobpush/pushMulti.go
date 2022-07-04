package mob_push_sdk



type PushMulti struct {
	PushMulti   []PushWork `json:"pushWork"`
}

type PushWork struct {
	WorkNo string `json:"workno"`
	Source string `json:"source"`
	PushTarget PushTarget `json:"pushTarget"`
	PushNotify PushNotify `json:"pushNotify"`
	PushOperator PushOperator `json:"pushOperator"`
	PushForward PushForward `json:"pushForward"`
}

