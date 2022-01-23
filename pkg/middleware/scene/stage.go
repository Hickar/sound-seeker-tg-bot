package scene

import (
	"fmt"

	sessionMiddleware "github.com/Hickar/sound-seeker-bot/pkg/middleware/session"
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
			sessionScene := stage.GetScene(ctx)
			if sessionScene == nil {
				sessionScene = newDefaultScene(stage)
				stage.SetScene(ctx, sessionScene.Name)
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

func (st *Stage) GetScene(ctx telebot.Context) *Scene {
	session := sessionMiddleware.GetSession(ctx)
	sceneName, ok := session.Get(sceneSessionKey).(string)
	if !ok {
		return nil
	}

	scene, ok := st.Scenes[sceneName]
	if !ok {
		return nil
	}

	return scene
}

func (st *Stage) EnterScene(ctx telebot.Context, sceneName string) error {
	session := sessionMiddleware.GetSession(ctx)
	newScene, ok := st.Scenes[sceneName]
	if !ok {
		return fmt.Errorf("no scene with name \"%s\"", sceneName)
	}

	session.Set(sceneSessionKey, sceneName)
	if newScene.startHandler != nil {
		return newScene.startHandler(ctx)
	}

	return nil
}

func (st *Stage) SetScene(ctx telebot.Context, sceneName string) {
	session := sessionMiddleware.GetSession(ctx)
	session.Set(sceneSessionKey, sceneName)
}
