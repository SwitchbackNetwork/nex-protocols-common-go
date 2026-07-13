package ticket_granting

import (
	"encoding/hex"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	ticket_granting "github.com/PretendoNetwork/nex-protocols-go/v2/ticket-granting"
	ticket_granting_types "github.com/PretendoNetwork/nex-protocols-go/v2/ticket-granting/types"
)

func (commonProtocol *CommonProtocol) validateAndRequestTicketWithParam(err error, packet nex.PacketInterface, callID uint32, param ticket_granting_types.ValidateAndRequestTicketParam) (*nex.RMCMessage, *nex.Error) {
	if commonProtocol.ValidateLoginData == nil {
		common_globals.Logger.Error("TicketGranting::ValidateAndRequestTicketWithParam missing ValidateLoginData!")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "TicketGranting::ValidateAndRequestTicketWithParam missing ValidateLoginData!")
	}

	if err != nil {
		common_globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.Core.InvalidArgument, err.Error())
	}

	connection := packet.Sender().(*nex.PRUDPConnection)
	endpoint := connection.Endpoint().(*nex.PRUDPEndPoint)

	sourceAccount, errorCode := endpoint.AccountDetailsByUsername(string(param.Username))

	if errorCode == nil {
		errorCode = commonProtocol.ValidateLoginData(sourceAccount.PID, param.ExtraData)
	}

	var targetAccount *nex.Account
	if errorCode == nil {
		targetAccount, errorCode = endpoint.AccountDetailsByUsername(commonProtocol.SecureServerAccount.Username)
	}

	var sourceKey []byte
	if errorCode == nil && sourceAccount.RequiresTokenAuth {
		sourceKey, errorCode = commonProtocol.SourceKeyFromToken(sourceAccount, param.ExtraData)
	}

	var encryptedTicket []byte
	if errorCode == nil {
		encryptedTicket, errorCode = generateTicket(sourceAccount, targetAccount, sourceKey, commonProtocol.SessionKeyLength, endpoint)
	}

	if errorCode != nil {
		common_globals.Logger.Error(errorCode.Message)
		return nil, nex.NewError(errorCode.ResultCode, errorCode.Message)
	}

	result := ticket_granting_types.NewValidateAndRequestTicketResult()
	result.SourcePID = sourceAccount.PID
	result.BufResponse = types.NewBuffer(encryptedTicket)
	result.ServiceNodeURL = commonProtocol.SecureStationURL.Copy().(types.StationURL)
	result.CurrentUTCTime = types.NewDateTime(0).Now()
	result.ReturnMsg = commonProtocol.BuildName.Copy().(types.String)

	if sourceKey != nil {
		result.SourceKey = types.String(hex.EncodeToString(sourceKey))
	}

	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	result.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = ticket_granting.ProtocolID
	rmcResponse.MethodID = ticket_granting.MethodLoginWithContext
	rmcResponse.CallID = callID

	if commonProtocol.OnAfterValidateAndRequestTicketWithParam != nil {
		go commonProtocol.OnAfterValidateAndRequestTicketWithParam(packet, param)
	}

	return rmcResponse, nil
}
