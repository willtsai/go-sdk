package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testActorType = "test"

func TestInvokeActor(t *testing.T) {
	ctx := context.Background()
	in := &InvokeActorRequest{
		ActorID:   "fn",
		Method:    "mockMethod",
		Data:      []byte(`{hello}`),
		ActorType: testActorType,
	}

	t.Run("invoke actor without data ", func(t *testing.T) {
		in.Data = nil
		out, err := testClient.InvokeActor(ctx, in)
		in.Data = []byte(`{hello}`)
		assert.Nil(t, err)
		assert.NotNil(t, out)
	})

	t.Run("invoke actor without method", func(t *testing.T) {
		in.Method = ""
		out, err := testClient.InvokeActor(ctx, in)
		in.Method = "mockMethod"
		assert.NotNil(t, err)
		assert.Nil(t, out)
	})

	t.Run("invoke actor without id ", func(t *testing.T) {
		in.ActorID = ""
		out, err := testClient.InvokeActor(ctx, in)
		in.ActorID = "fn"
		assert.NotNil(t, err)
		assert.Nil(t, out)
	})

	t.Run("invoke actor without type", func(t *testing.T) {
		in.ActorType = ""
		out, err := testClient.InvokeActor(ctx, in)
		in.ActorType = testActorType
		assert.NotNil(t, err)
		assert.Nil(t, out)
	})

	t.Run("invoke actor without empty input", func(t *testing.T) {
		in = nil
		out, err := testClient.InvokeActor(ctx, in)
		assert.NotNil(t, err)
		assert.Nil(t, out)
	})
}

func TestRegisterActorReminder(t *testing.T) {
	ctx := context.Background()
	in := &RegisterActorReminderRequest{
		ActorID:   "fn",
		Data:      []byte(`{hello}`),
		ActorType: testActorType,
		Name:      "mockName",
		Period:    "2s",
		DueTime:   "4s",
	}

	t.Run("invoke register actor reminder without actorType", func(t *testing.T) {
		in.ActorType = ""
		err := testClient.RegisterActorReminder(ctx, in)
		in.ActorType = testActorType
		assert.NotNil(t, err)
	})

	t.Run("invoke register actor reminder without id ", func(t *testing.T) {
		in.ActorID = ""
		err := testClient.RegisterActorReminder(ctx, in)
		in.ActorID = "fn"
		assert.NotNil(t, err)
	})

	t.Run("invoke register actor reminder without Name ", func(t *testing.T) {
		in.Name = ""
		err := testClient.RegisterActorReminder(ctx, in)
		in.Name = "mockName"
		assert.NotNil(t, err)
	})

	t.Run("invoke register actor reminder without period ", func(t *testing.T) {
		in.ActorType = ""
		err := testClient.RegisterActorReminder(ctx, in)
		in.ActorType = testActorType
		assert.NotNil(t, err)
	})

	t.Run("invoke register actor reminder without dutTime ", func(t *testing.T) {
		in.DueTime = ""
		err := testClient.RegisterActorReminder(ctx, in)
		in.DueTime = "2s"
		assert.NotNil(t, err)
	})

	t.Run("invoke register actor reminder ", func(t *testing.T) {
		assert.Nil(t, testClient.RegisterActorReminder(ctx, in))
	})

	t.Run("invoke register actor reminder with empty param", func(t *testing.T) {
		assert.NotNil(t, testClient.RegisterActorReminder(ctx, nil))
	})
}

func TestRegisterActorTimer(t *testing.T) {
	ctx := context.Background()
	in := &RegisterActorTimerRequest{
		ActorID:   "fn",
		Data:      []byte(`{hello}`),
		ActorType: testActorType,
		Name:      "mockName",
		Period:    "2s",
		DueTime:   "4s",
		CallBack:  "mockFunc",
	}

	t.Run("invoke register actor timer without actorType", func(t *testing.T) {
		in.ActorType = ""
		err := testClient.RegisterActorTimer(ctx, in)
		in.ActorType = testActorType
		assert.NotNil(t, err)
	})

	t.Run("invoke register actor timer without id ", func(t *testing.T) {
		in.ActorID = ""
		err := testClient.RegisterActorTimer(ctx, in)
		in.ActorID = "fn"
		assert.NotNil(t, err)
	})

	t.Run("invoke register actor timer without Name ", func(t *testing.T) {
		in.Name = ""
		err := testClient.RegisterActorTimer(ctx, in)
		in.Name = "mockName"
		assert.NotNil(t, err)
	})

	t.Run("invoke register actor timer without period ", func(t *testing.T) {
		in.ActorType = ""
		err := testClient.RegisterActorTimer(ctx, in)
		in.ActorType = testActorType
		assert.NotNil(t, err)
	})

	t.Run("invoke register actor timer without dutTime ", func(t *testing.T) {
		in.DueTime = ""
		err := testClient.RegisterActorTimer(ctx, in)
		in.DueTime = "2s"
		assert.NotNil(t, err)
	})

	t.Run("invoke register actor timer without callBack ", func(t *testing.T) {
		in.CallBack = ""
		err := testClient.RegisterActorTimer(ctx, in)
		in.CallBack = "mockFunc"
		assert.NotNil(t, err)
	})

	t.Run("invoke register actor timer without data ", func(t *testing.T) {
		in.Data = nil
		err := testClient.RegisterActorTimer(ctx, in)
		in.Data = []byte(`{hello}`)
		assert.Nil(t, err)
	})

	t.Run("invoke register actor timer", func(t *testing.T) {
		assert.Nil(t, testClient.RegisterActorTimer(ctx, in))
	})

	t.Run("invoke register actor timer with empty param", func(t *testing.T) {
		assert.NotNil(t, testClient.RegisterActorTimer(ctx, nil))
	})
}

func TestUnregisterActorReminder(t *testing.T) {
	ctx := context.Background()
	in := &UnregisterActorReminderRequest{
		ActorID:   "fn",
		ActorType: testActorType,
		Name:      "mockName",
	}

	t.Run("invoke unregister actor reminder without actorType", func(t *testing.T) {
		in.ActorType = ""
		err := testClient.UnregisterActorReminder(ctx, in)
		in.ActorType = testActorType
		assert.NotNil(t, err)
	})

	t.Run("invoke unregister actor reminder without id ", func(t *testing.T) {
		in.ActorID = ""
		err := testClient.UnregisterActorReminder(ctx, in)
		in.ActorID = "fn"
		assert.NotNil(t, err)
	})

	t.Run("invoke unregister actor reminder without Name ", func(t *testing.T) {
		in.Name = ""
		err := testClient.UnregisterActorReminder(ctx, in)
		in.Name = "mockName"
		assert.NotNil(t, err)
	})

	t.Run("invoke unregister actor reminder without period ", func(t *testing.T) {
		in.ActorType = ""
		err := testClient.UnregisterActorReminder(ctx, in)
		in.ActorType = testActorType
		assert.NotNil(t, err)
	})

	t.Run("invoke unregister actor reminder ", func(t *testing.T) {
		assert.Nil(t, testClient.UnregisterActorReminder(ctx, in))
	})

	t.Run("invoke unregister actor reminder with empty param", func(t *testing.T) {
		assert.NotNil(t, testClient.UnregisterActorReminder(ctx, nil))
	})
}

func TestUnregisterActorTimer(t *testing.T) {
	ctx := context.Background()
	in := &UnregisterActorTimerRequest{
		ActorID:   "fn",
		ActorType: testActorType,
		Name:      "mockName",
	}

	t.Run("invoke unregister actor timer without actorType", func(t *testing.T) {
		in.ActorType = ""
		err := testClient.UnregisterActorTimer(ctx, in)
		in.ActorType = testActorType
		assert.NotNil(t, err)
	})

	t.Run("invoke register actor timer without id ", func(t *testing.T) {
		in.ActorID = ""
		err := testClient.UnregisterActorTimer(ctx, in)
		in.ActorID = "fn"
		assert.NotNil(t, err)
	})

	t.Run("invoke register actor timer without Name ", func(t *testing.T) {
		in.Name = ""
		err := testClient.UnregisterActorTimer(ctx, in)
		in.Name = "mockName"
		assert.NotNil(t, err)
	})

	t.Run("invoke register actor timer without period ", func(t *testing.T) {
		in.ActorType = ""
		err := testClient.UnregisterActorTimer(ctx, in)
		in.ActorType = testActorType
		assert.NotNil(t, err)
	})

	t.Run("invoke register actor timer ", func(t *testing.T) {
		assert.Nil(t, testClient.UnregisterActorTimer(ctx, in))
	})

	t.Run("invoke register actor timer with empty param", func(t *testing.T) {
		assert.NotNil(t, testClient.UnregisterActorTimer(ctx, nil))
	})
}
