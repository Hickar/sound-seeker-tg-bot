package scene

import (
	"gopkg.in/tucnak/telebot.v3"
)

const (
	sceneSessionKey = "scene_session"
	defaultScene    = "scene_default"
)

type Stage struct {
	Scenes map[string]*Scene
}

func NewStage() *Stage {
	return &Stage{Scenes: map[string]*Scene{}}
}

func (st *Stage) Register(scenes ...*Scene) {
	for _, scene := range scenes {
		scene.Stage = st
		st.Scenes[scene.Name] = scene
	}
}

func Middleware(stage *Stage) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(ctx telebot.Context) error {
			sessionScene := GetScene(ctx)
			if sessionScene == nil {
				sessionScene = newDefaultScene(stage)
				SetScene(ctx, sessionScene)
			}

			currentScene, ok := stage.Scenes[sessionScene.Name]
			if sessionScene.Name == defaultScene || !ok {
				return next(ctx)
			}

			currentScene.ProcessUpdate(ctx.Update())

			return nil
		}
	}
}