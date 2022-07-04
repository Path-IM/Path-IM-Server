package mob_push_sdk

type requestData map[string]interface{}

func (client *PushClient) NewRequestData()  requestData{
	requestData := make(requestData)
	requestData["appkey"] = client.AppKey
	return requestData
}
