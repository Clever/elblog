package schemas

import (
	"github.com/Clever/elblog"
	"strconv"
	"strings"
	"time"
)

// ALBLogSchema is a representation of a row in access-logs-alb-global for use with Parquet:
// https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-access-logs.html
// https://github.com/xitongsys/parquet-go
type ALBLogSchema struct {
	Type                   string  `parquet:"name=type, type=UTF8, encoding=PLAIN_DICTIONARY"`
	Time                   string  `parquet:"name=time, type=UTF8"`
	ELB                    string  `parquet:"name=elb, type=UTF8, encoding=PLAIN_DICTIONARY"`
	ClientIP               string  `parquet:"name=client_ip, type=UTF8"`
	ClientPort             int32   `parquet:"name=client_port, type=INT32"`
	TargetIP               string  `parquet:"name=target_ip, type=UTF8"`
	TargetPort             int32   `parquet:"name=target_port, type=INT32"`
	RequestProcessingTime  float64 `parquet:"name=request_processing_time, type=DOUBLE"`
	TargetProcessingTime   float64 `parquet:"name=target_processing_time, type=DOUBLE"`
	ResponseProcessingTime float64 `parquet:"name=response_processing_time, type=DOUBLE"`
	ELBStatusCode          string  `parquet:"name=elb_status_code, type=UTF8, encoding=PLAIN_DICTIONARY"`
	TargetStatusCode       string  `parquet:"name=target_status_code, type=UTF8, encoding=PLAIN_DICTIONARY"`
	ReceivedBytes          int64   `parquet:"name=received_bytes, type=INT64"`
	SentBytes              int64   `parquet:"name=sent_bytes, type=INT64"`
	RequestVerb            string  `parquet:"name=request_verb, type=UTF8, encoding=PLAIN_DICTIONARY"`
	RequestURL             string  `parquet:"name=request_url, type=UTF8"`
	RequestProto           string  `parquet:"name=request_proto, type=UTF8, encoding=PLAIN_DICTIONARY"`
	UserAgent              string  `parquet:"name=user_agent, type=UTF8"`
	SSLCipher              string  `parquet:"name=ssl_cipher, type=UTF8, encoding=PLAIN_DICTIONARY"`
	SSLProtocol            string  `parquet:"name=ssl_protocol, type=UTF8, encoding=PLAIN_DICTIONARY"`
	TargetGroupARN         string  `parquet:"name=target_group_arn, type=UTF8, encoding=PLAIN_DICTIONARY"`
	TraceID                string  `parquet:"name=trace_id, type=UTF8"`
	DomainName             string  `parquet:"name=domain_name, type=UTF8, encoding=PLAIN_DICTIONARY"`
	ChosenCertARN          string  `parquet:"name=chosen_cert_arn, type=UTF8"`
	MatchedRulePriority    string  `parquet:"name=matched_rule_priority, type=UTF8, encoding=PLAIN_DICTIONARY"`
	RequestCreationTime    string  `parquet:"name=request_creation_time, type=UTF8"`
	ActionsExecuted        string  `parquet:"name=actions_executed, type=UTF8, encoding=PLAIN_DICTIONARY"`
	RedirectURL            string  `parquet:"name=redirect_url, type=UTF8"`
	ErrorReason            string  `parquet:"name=error_reason, type=UTF8, encoding=PLAIN_DICTIONARY"`
	OtherFields            string  `parquet:"name=other_fields, type=UTF8"`
}

// ELBLogToALBLogSchema converts an elblog to an ALBLogSchema that has tags for parquet
func ELBLogToALBLogSchema(log elblog.Log) ALBLogSchema {
	return ALBLogSchema{
		Type:                   log.Type,
		Time:                   log.Time.Format(time.RFC3339),
		ELB:                    log.Name,
		ClientIP:               log.From.IP.String(),
		ClientPort:             int32(log.From.Port),
		TargetIP:               log.To.IP.String(),
		TargetPort:             int32(log.To.Port),
		RequestProcessingTime:  log.RequestProcessingTime.Seconds(),
		TargetProcessingTime:   log.BackendProcessingTime.Seconds(),
		ResponseProcessingTime: log.ResponseProcessingTime.Seconds(),
		ELBStatusCode:          strconv.Itoa(log.ELBStatusCode),
		TargetStatusCode:       log.BackendStatusCode,
		ReceivedBytes:          log.ReceivedBytes,
		SentBytes:              log.SentBytes,
		RequestVerb:            strings.Split(log.Request, " ")[0],
		RequestURL:             strings.Split(log.Request, " ")[1],
		RequestProto:           strings.Split(log.Request, " ")[2],
		UserAgent:              log.UserAgent,
		SSLCipher:              log.SSLCipher,
		SSLProtocol:            log.SSLProtocol,
		TargetGroupARN:         log.TargetGroupARN,
		TraceID:                log.TraceID,
		DomainName:             log.DomainName,
		ChosenCertARN:          log.ChosenCertARN,
		MatchedRulePriority:    log.MatchedRulePriority,
		RequestCreationTime:    log.RequestCreationTime,
		ActionsExecuted:        log.ActionsExecuted,
		RedirectURL:            log.RedirectURL,
		ErrorReason:            log.ErrorReason,
		OtherFields:            log.OtherFields,
	}
}
