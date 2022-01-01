package main

import (
	"fmt"

	math "github.com/chewxy/math32"
	RLBot "github.com/Trey2k/RLBotGo"
)

var lastTouch float32
var totalTouches int = 0

// getInput takes in a GameState which contains the gameTickPacket, ballPredidctions, fieldInfo and matchSettings
// it also takes in the RLBot object. And returns a PlayerInput

func getInput(gameState *RLBot.GameState, rlBot *RLBot.RLBot) *RLBot.ControllerState {
	PlayerInput := &RLBot.ControllerState{}

	selfLocation := gameState.GameTick.Players[rlBot.PlayerIndex].Physics.Location
	selfRotation := gameState.GameTick.Players[rlBot.PlayerIndex].Physics.Rotation
	ballLocation := gameState.GameTick.Ball.Physics.Location

	if !gameState.GameTick.Players[rlBot.PlayerIndex].HasWheelContact {
		PlayerInput.Roll = -1.0*selfRotation.Roll/math.Pi
	}

	localX := ballLocation.X - selfLocation.X
	localY := ballLocation.Y - selfLocation.Y
	toBallAngle := math.Atan2(localY,localX)
	steer := toBallAngle - selfRotation.Yaw

	if steer < -math.Pi {
		steer += math.Pi * 2.0;
	} else if steer >= math.Pi {
		steer -= math.Pi * 2.0;
	}

	if (steer > 1) {
		steer = 1
	} else if (steer < -1) {
		steer = -1
	}

	PlayerInput.Steer = steer

	distanceToBall := math.Sqrt(localX * localX + localY * localY)

	if distanceToBall < 300 && math.Abs(steer) > 0.3 {
		PlayerInput.Handbrake = true
	} 

	PlayerInput.Throttle = 1.0

	return PlayerInput

}

func main() {

	// connect to RLBot
	rlBot, err := RLBot.Connect(23234)
	if err != nil {
		panic(err)
	}

	// Send ready message
	err = rlBot.SendReadyMessage(true, true, true)
	if err != nil {
		panic(err)
	}

	// Set our tick handler
	err = rlBot.SetGetInput(getInput)
	fmt.Println(err.Error())

}
