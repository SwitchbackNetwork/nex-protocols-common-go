package message_delivery

import (
	"github.com/PretendoNetwork/nex-go/v2"
	messaging "github.com/PretendoNetwork/nex-protocols-go/v2/messaging"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	messaging_database "github.com/PretendoNetwork/nex-protocols-common-go/v2/messaging/database"
)

type CommonProtocol struct {
	endpoint              nex.EndpointInterface
	protocol              messaging.Interface
	manager               *common_globals.MessagingManager
}

// SetManager defines the messaging manager to be used by the common protocol
func (commonProtocol *CommonProtocol) SetManager(manager *common_globals.MessagingManager) {
	var err error

	commonProtocol.manager = manager

	if manager.ProcessMessage == nil {
		manager.ProcessMessage = messaging_database.ProcessMessage
	}

	_, err = manager.Database.Exec(`CREATE SCHEMA IF NOT EXISTS messaging`)
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return
	}

	_, err = manager.Database.Exec(`CREATE TABLE IF NOT EXISTS messaging.instant_messages (
		id bigserial PRIMARY KEY,
		recipient_id numeric(20),
		recipient_type numeric(10),
		parent_id bigint,
		sender_pid numeric(20),
		reception_time timestamp,
		lifetime numeric(10),
		flags bigint,
		subject text,
		sender text,
		type text NOT NULL DEFAULT ''
	)`)
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return
	}

	_, err = manager.Database.Exec(`CREATE TABLE IF NOT EXISTS messaging.instant_text_messages (
		id bigint PRIMARY KEY REFERENCES messaging.instant_messages(id),
		body text
	)`)
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return
	}

	_, err = manager.Database.Exec(`CREATE TABLE IF NOT EXISTS messaging.instant_binary_messages (
		id bigint PRIMARY KEY REFERENCES messaging.instant_messages(id),
		body bytea
	)`)
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return
	}

	_, err = manager.Database.Exec(`CREATE TABLE IF NOT EXISTS messaging.messages (
		id bigserial PRIMARY KEY,
		recipient_id numeric(20),
		recipient_type numeric(10),
		parent_id bigint,
		sender_pid numeric(20),
		reception_time timestamp,
		lifetime numeric(10),
		flags bigint,
		subject text,
		sender text,
		type text NOT NULL DEFAULT '',
		read boolean NOT NULL DEFAULT false,
		deleted boolean NOT NULL DEFAULT false
	)`)
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return
	}

	_, err = manager.Database.Exec(`CREATE TABLE IF NOT EXISTS messaging.text_messages (
		id bigint PRIMARY KEY REFERENCES messaging.messages(id),
		body text
	)`)
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return
	}

	_, err = manager.Database.Exec(`CREATE TABLE IF NOT EXISTS messaging.binary_messages (
		id bigint PRIMARY KEY REFERENCES messaging.messages(id),
		body bytea
	)`)
	if err != nil {
		common_globals.Logger.Error(err.Error())
		return
	}
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol messaging.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerDeliverMessage(commonProtocol.deliverMessage)
	protocol.SetHandlerDeliverMessageMultiTarget(commonProtocol.deliverMessageMultiTarget)

	return commonProtocol
}
