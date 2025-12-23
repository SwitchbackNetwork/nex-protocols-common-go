package match_making_ext

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	match_making_ext "github.com/PretendoNetwork/nex-protocols-go/v2/match-making-ext"
	common_globals "github.com/SwitchbackNetwork/nex-protocols-common-go/v2/globals"
)

type CommonProtocol struct {
	endpoint                nex.EndpointInterface
	protocol                match_making_ext.Interface
	manager                 *common_globals.MatchmakingManager
	OnAfterEndParticipation func(acket nex.PacketInterface, idGathering types.UInt32, strMessage types.String)
}

// SetManager defines the matchmaking manager to be used by the common protocol
func (commonProtocol *CommonProtocol) SetManager(manager *common_globals.MatchmakingManager) {
	commonProtocol.manager = manager
}

// NewCommonProtocol returns a new CommonProtocol
func NewCommonProtocol(protocol match_making_ext.Interface) *CommonProtocol {
	commonProtocol := &CommonProtocol{
		endpoint: protocol.Endpoint(),
		protocol: protocol,
	}

	protocol.SetHandlerEndParticipation(commonProtocol.endParticipation)

	return commonProtocol
}
