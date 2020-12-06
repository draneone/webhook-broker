package data

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	sampleChannel             = &Channel{}
	sampleCallbackURL         = getSampleURL("http://imytech.net/")
	sampleRelativeCallbackURL = getSampleURL("./")
	getSampleURL              = func(sampleURL string) *url.URL {
		url, _ := url.Parse(sampleURL)
		return url
	}
)

func TestNewConsumer(t *testing.T) {
	t.Run("EmptyID", func(t *testing.T) {
		t.Parallel()
		_, err := NewConsumer(sampleChannel, "", "", sampleCallbackURL)
		assert.Equal(t, ErrInsufficientInformationForCreating, err)
	})
	t.Run("EmptyToken", func(t *testing.T) {
		t.Parallel()
		_, err := NewConsumer(sampleChannel, someID, "", sampleCallbackURL)
		assert.Equal(t, ErrInsufficientInformationForCreating, err)
	})
	t.Run("NilChannel", func(t *testing.T) {
		t.Parallel()
		_, err := NewConsumer(nil, someID, someToken, sampleCallbackURL)
		assert.Equal(t, ErrInsufficientInformationForCreating, err)
	})
	t.Run("RelativeURL", func(t *testing.T) {
		t.Parallel()
		_, err := NewConsumer(sampleChannel, someID, someToken, sampleRelativeCallbackURL)
		assert.Equal(t, ErrInsufficientInformationForCreating, err)
	})
	t.Run("Valid", func(t *testing.T) {
		t.Parallel()
		consumer, err := NewConsumer(sampleChannel, someID, someToken, sampleCallbackURL)
		assert.Nil(t, err)
		assert.NotNil(t, consumer.ID)
		assert.Equal(t, someID, consumer.ConsumerID)
		assert.Equal(t, someID, consumer.Name)
		assert.Equal(t, someToken, consumer.Token)
	})
}

func TestConsumerIsInValidState(t *testing.T) {
	t.Run("True", func(t *testing.T) {
		t.Parallel()
		consumer, _ := NewConsumer(sampleChannel, someID, someToken, sampleCallbackURL)
		assert.True(t, consumer.IsInValidState())
	})
	t.Run("EmptyIDFalse", func(t *testing.T) {
		t.Parallel()
		consumer, _ := NewConsumer(sampleChannel, someID, someToken, sampleCallbackURL)
		consumer.ConsumerID = ""
		assert.False(t, consumer.IsInValidState())
	})
	t.Run("NilChannelFalse", func(t *testing.T) {
		t.Parallel()
		consumer, _ := NewConsumer(sampleChannel, someID, someToken, sampleCallbackURL)
		consumer.ConsumingFrom = nil
		assert.False(t, consumer.IsInValidState())
	})
	t.Run("RelativeURLFalse", func(t *testing.T) {
		t.Parallel()
		consumer, _ := NewConsumer(sampleChannel, someID, someToken, sampleCallbackURL)
		consumer.CallbackURL = sampleRelativeCallbackURL.String()
		assert.False(t, consumer.IsInValidState())
	})
}

func TestConsumerQuickFix(t *testing.T) {
	t.Parallel()
	consumer, _ := NewConsumer(sampleChannel, someID, someToken, sampleCallbackURL)
	consumer.Name = ""
	assert.False(t, consumer.IsInValidState())
	assert.True(t, len(consumer.Name) <= 0)
	assert.True(t, consumer.QuickFix())
	assert.True(t, consumer.IsInValidState())
	assert.Equal(t, someID, consumer.Name)
}