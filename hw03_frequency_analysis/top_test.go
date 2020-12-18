package hw03_frequency_analysis //nolint:golint

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		top, err := Top10("", false)
		assert.NoError(t, err)
		assert.Len(t, top, 0)
	})

	t.Run("positive test w/o asterisk", func(t *testing.T) {
		expected := []string{"он", "и", "а", "что", "ты", "не", "если", "-", "то", "Кристофер"}
		top, err := Top10(text, false)
		assert.NoError(t, err)
		assert.ElementsMatch(t, expected, top)
	})

	t.Run("positive test w/ asterisk", func(t *testing.T) {
		expected := []string{"он", "а", "и", "что", "ты", "не", "если", "то", "его", "кристофер", "робин", "в"}
		top, err := Top10(text, true)
		assert.NoError(t, err)
		assert.Subset(t, expected, top)
	})

	t.Run("less than 10", func(t *testing.T) {
		input := "a b c"
		want := []string{"a", "b", "c"}
		got, err := Top10(input, true)
		assert.NoError(t, err)
		assert.Subset(t, got, want)
	})

	t.Run("deutcsh", func(t *testing.T) {
		input := "Schulabgänger? Schulabgänger? Student? Absolvent? Berufserfahrener?"
		want := []string{"schulabgänger", "absolvent", "berufserfahrener", "student"}
		got, err := Top10(input, true)
		assert.NoError(t, err)
		assert.Subset(t, got, want)
	})
}

func TestNewParser(t *testing.T) {
	t.Run("wrong pattern returns error", func(t *testing.T) {
		_, err := newParser(`[\p{L}\d`, false)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrRegExCompilation))
	})
}
