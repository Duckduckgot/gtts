package htgotts

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/Duckduckgot/gtts/handlers"
)

/**
 * Required:
 * - mplayer
 *
 * Use:
 *
 * speech := htgotts.Speech{Folder: "audio", Language: "en", Handler: MPlayer}
 */

// Speech struct
type Speech struct {
	Folder   string
	Language string
	Handler  handlers.PlayerInterface
}

// Creates a speech file with a given name
func (speech *Speech) CreateSpeechFile(text string, fileName string) (string, error) {
	var err error

	f := speech.Folder + "/" + fileName + ".mp3"
	if err = speech.createFolderIfNotExists(speech.Folder); err != nil {
		return "", err
	}

	if err = speech.downloadIfNotExists(f, text); err != nil {
		return "", err
	}

	return f, nil
}

// Plays an existent .mp3 file
func (speech *Speech) PlaySpeechFile(fileName string) error {
	if speech.Handler == nil {
		mplayer := handlers.MPlayer{}
		return mplayer.Play(fileName)
	}

	return speech.Handler.Play(fileName)
}

// Speak downloads speech and plays it using mplayer
func (speech *Speech) Speak(text string) error {

	var err error
	generatedHashName := speech.generateHashName(text)

	fileName, err := speech.CreateSpeechFile(text, generatedHashName)
	if err != nil {
		return err
	}

	return speech.PlaySpeechFile(fileName)
}

/**
 * Create the folder if does not exists.
 */
func (speech *Speech) createFolderIfNotExists(folder string) error {
	dir, err := os.Open(folder)
	if os.IsNotExist(err) {
		return os.MkdirAll(folder, 0700)
	}

	dir.Close()
	return nil
}

/**
 * Download the voice file if does not exists.
 */
func (speech *Speech) downloadIfNotExists(fileName string, text string) error {
	f, err := os.Open(fileName)
	if err != nil {
		url := fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), speech.Language)
		response, err := http.Get(url)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		output, err := os.Create(fileName)
		if err != nil {
			return err
		}

		_, err = io.Copy(output, response.Body)
		return err
	}

	f.Close()
	return nil
}

// Speak downloads speech and plays it using mplayer
func (speech *Speech) SpeakB(text string) ([]byte, error) {

	var err error

	speakbyte, err := speech.downloadByte(text)
	if err != nil {
		return nil, err
	}

	return speakbyte, nil
}

/**
 * Download the voice as byte.
 */
func (speech *Speech) downloadByte(text string) ([]byte, error) {
	url := fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), speech.Language)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (speech *Speech) generateHashName(name string) string {
	hash := md5.Sum([]byte(name))
	return fmt.Sprintf("%s_%s", speech.Language, hex.EncodeToString(hash[:]))
}
