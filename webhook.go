package paypal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	kWebHookAPI                = "/v1/notifications/webhooks"
	kVerityWebHookSignatureAPI = "/v1/notifications/verify-webhook-signature"
)

// CreateWebhook https://developer.paypal.com/docs/api/webhooks/#webhooks_create
func (this *Client) CreateWebhook(callBackURL string, eventTypeList ...string) (results *Webhook, err error) {
	var api = this.BuildAPI(kWebHookAPI)
	var p = &Webhook{}

	var events = make([]*EventType, 0, len(eventTypeList))
	for _, name := range eventTypeList {
		var event = &EventType{}
		event.Name = name
		events = append(events, event)
	}
	p.EventTypes = events

	p.URL = callBackURL
	err = this.doRequestWithAuth(http.MethodPost, api, p, &results)
	return results, err
}

// GetWebhookList https://developer.paypal.com/docs/api/webhooks/#webhooks_get-all
func (this *Client) GetWebhookList() (results *WebhookList, err error) {
	var api = this.BuildAPI(kWebHookAPI)
	err = this.doRequestWithAuth(http.MethodGet, api, nil, &results)
	return results, err
}

// GetWebhook https://developer.paypal.com/docs/api/webhooks/#webhooks_get
func (this *Client) GetWebhookDetails(webhookId string) (results *Webhook, err error) {
	var api = this.BuildAPI(kWebHookAPI, webhookId)
	err = this.doRequestWithAuth(http.MethodGet, api, nil, &results)
	return results, err
}

// DeleteWebhook https://developer.paypal.com/docs/api/webhooks/#webhooks_delete
func (this *Client) DeleteWebhook(webhookId string) (err error) {
	var api = this.BuildAPI(kWebHookAPI, webhookId)
	err = this.doRequestWithAuth(http.MethodDelete, api, nil, nil)
	return err
}

// verifyWebhookSignature https://developer.paypal.com/docs/api/webhooks/#verify-webhook-signature_post
func (this *Client) verifyWebhookSignature(param *verifyWebhookSignatureParam) (results *verifyWebhookSignatureResponse, err error) {
	var api = this.BuildAPI(kVerityWebHookSignatureAPI)
	err = this.doRequestWithAuth(http.MethodPost, api, param, &results)
	return results, err
}

// GetWebhookEvent 用于处理 webbook 回调
func (this *Client) GetWebhookEvent(webhookId string, req *http.Request) (event *Event, err error) {
	req.ParseForm()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		return nil, err
	}

	var rawRes json.RawMessage
	event = &Event{
		Resource: &rawRes,
	}

	if err = json.Unmarshal(body, &event); err != nil {
		return nil, err
	}

	if event == nil || (event != nil && (event.Id == "" || event.EventType == "")) {
		return nil, errors.New("unknown webhook event")
	}

	switch event.ResourceType {
	case K_EVENT_RESOURCE_TYPE_SALE:
		var sale *Sale
		if err = json.Unmarshal(rawRes, &sale); err != nil {
			return nil, err
		}
		event.Resource = sale
	case K_EVENT_RESOURCE_TYPE_REFUND:
		var refund *Refund
		if err = json.Unmarshal(rawRes, &refund); err != nil {
			return nil, err
		}
		event.Resource = refund
	case K_EVENT_RESOURCE_TYPE_INVOICES:
		var invoice *Invoice
		if err = json.Unmarshal(rawRes, &invoice); err != nil {
			return nil, err
		}
		event.Resource = invoice
	case K_EVENT_RESOURCE_TYPE_DISPUTE:
		var dispute *Dispute
		if err = json.Unmarshal(rawRes, &dispute); err != nil {
			return nil, err
		}
		event.Resource = dispute
	default:
		var data map[string]interface{}
		if err = json.Unmarshal(rawRes, &data); err != nil {
			return nil, err
		}
		event.Resource = data
	}

	var verifyParam = &verifyWebhookSignatureParam{}
	verifyParam.AuthAlgo = req.Header.Get("Paypal-Auth-Algo")
	verifyParam.CertURL = req.Header.Get("Paypal-Cert-Url")
	verifyParam.TransmissionId = req.Header.Get("Paypal-Transmission-Id")
	verifyParam.TransmissionSig = req.Header.Get("Paypal-Transmission-Sig")
	verifyParam.TransmissionTime = req.Header.Get("Paypal-Transmission-Time")
	verifyParam.WebhookId = webhookId
	verifyParam.WebhookEvent = jsonString(string(body))

	verifyResp, err := this.verifyWebhookSignature(verifyParam)
	if err != nil {
		return nil, err
	}

	if verifyResp.VerificationStatus != kSuccess {
		return nil, errors.New(fmt.Sprintf("verify webhook %s", verifyResp.VerificationStatus))
	}

	return event, nil
}
