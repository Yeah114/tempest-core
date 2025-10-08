package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Yeah114/FunShuttler/core/minecraft/protocol"
	fpacket "github.com/Yeah114/FunShuttler/core/minecraft/protocol/packet"
	"github.com/Yeah114/FunShuttler/game_control/resources_control"
	uqdefines "github.com/Yeah114/FunShuttler/uqholder/defines"
	"github.com/Yeah114/tempest-core/network/app"
)

const abilityMask = protocol.AbilityBuild |
	protocol.AbilityMine |
	protocol.AbilityDoorsAndSwitches |
	protocol.AbilityOpenContainers |
	protocol.AbilityAttackPlayers |
	protocol.AbilityAttackMobs |
	protocol.AbilityOperatorCommands |
	protocol.AbilityTeleport

func fetchPlayerByUUID(state *app.FatalderState, uuidStr string) (uqdefines.PlayerUQReader, error) {
	var player uqdefines.PlayerUQReader
	err := state.WithResources(func(res *resources_control.Resources) error {
		holder := res.UQHolder()
		if holder == nil {
			return errors.New("uqholder unavailable")
		}
		micro := holder.Micro()
		if micro == nil {
			return errors.New("micro uqholder unavailable")
		}
		p, found := micro.GetPlayersInfo().GetPlayerByUUIDString(uuidStr)
		if !found {
			return wrapNotFound("player not found")
		}
		player = p
		return nil
	})
	return player, err
}

func fetchPlayerByName(state *app.FatalderState, name string) (uqdefines.PlayerUQReader, error) {
	var player uqdefines.PlayerUQReader
	err := state.WithResources(func(res *resources_control.Resources) error {
		holder := res.UQHolder()
		if holder == nil {
			return errors.New("uqholder unavailable")
		}
		micro := holder.Micro()
		if micro == nil {
			return errors.New("micro uqholder unavailable")
		}
		p, found := micro.GetPlayersInfo().GetPlayerByName(name)
		if !found {
			return wrapNotFound("player not found")
		}
		player = p
		return nil
	})
	return player, err
}

func abilityEnabled(player uqdefines.PlayerUQReader, flag uint32) (bool, error) {
	if player == nil {
		return false, errors.New("player not bound")
	}
	values, ok := player.GetValues()
	if !ok {
		return false, errors.New("ability values unavailable")
	}
	return (values & flag) != 0, nil
}

func updateAbility(state *app.FatalderState, player uqdefines.PlayerUQReader, flag uint32, allow bool) error {
	if player == nil {
		return errors.New("player not bound")
	}
	entityID, ok := player.GetEntityUniqueID()
	if !ok {
		return errors.New("entity unique id unavailable")
	}
	values, ok := player.GetValues()
	if !ok {
		return errors.New("ability values unavailable")
	}
	if allow {
		values |= flag
	} else {
		values &^= flag
	}
	requested := uint16(values & abilityMask)
	level := permissionLevelFor(values & abilityMask)

	return state.WithResources(func(res *resources_control.Resources) error {
		return res.WritePacket(&fpacket.RequestPermissions{
			EntityUniqueID:       entityID,
			PermissionLevel:      level,
			RequestedPermissions: requested,
		})
	})
}

func permissionLevelFor(mask uint32) uint8 {
	switch mask {
	case 0:
		return uint8(fpacket.PermissionLevelVisitor)
	case abilityMask:
		return uint8(fpacket.PermissionLevelOperator)
	default:
		return uint8(fpacket.PermissionLevelMember)
	}
}

func buildRawText(message string) (string, error) {
	nodes := map[string]any{
		"rawtext": []map[string]string{
			{"text": message},
		},
	}
	bs, err := json.Marshal(nodes)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func quotedCommandTarget(name string) string {
	name = strings.ReplaceAll(name, `"`, `\"`)
	return fmt.Sprintf(`"%s"`, name)
}
