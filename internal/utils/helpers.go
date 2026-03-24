package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/tekofx/crossposter/internal/config"
	merrors "github.com/tekofx/crossposter/internal/errors"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"

	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func SendMessageToOwner(ctx *th.Context, text string) *telego.Message {
	msg, err := ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(int64(config.Conf.TelegramOwner)),
		text,
	))

	if err != nil {
		logger.Fatal(merrors.TelegramCannotSendMessageToOwnerErrorCode, err.Error())
		return nil
	}
	return msg
}

func SendMessageToOwnerUsingBot(bot *telego.Bot, text string) *telego.Message {
	msg, err := bot.SendMessage(context.Background(), tu.Message(
		tu.ID(int64(config.Conf.TelegramOwner)),
		text,
	))

	if err != nil {
		logger.Fatal(merrors.New(merrors.TelegramCannotSendMessageToOwnerErrorCode, err.Error()))
		return nil
	}
	return msg
}

func SendPostToOwner(bot *telego.Bot, post *model.Post) *merrors.MError {
	if post.HasImages {
		var media []telego.InputMedia
		for _, image := range post.Images {
			downloadedImage, err := os.Open(image.Filename)
			if err != nil {
				return merrors.New(merrors.CannotReadFileErrorCode, err.Error())
			}
			defer downloadedImage.Close()

			inputFile := telego.InputFile{
				File: downloadedImage,
			}

			photo := &telego.InputMediaPhoto{
				Type:  "photo",
				Media: inputFile,
			}
			media = append(media, photo)
		}

		mediaGroup := telego.SendMediaGroupParams{
			ChatID: telego.ChatID{ID: int64(config.Conf.TelegramOwner)},
			Media:  media,
		}

		_, err := bot.SendMediaGroup(context.Background(), &mediaGroup)
		if err != nil {
			return merrors.New(merrors.TelegramCannotSendMediaGroupErrorCode, err.Error())
		}
	}

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("Editar").WithCallbackData("edit"),
			tu.InlineKeyboardButton("Borrar").WithCallbackData(fmt.Sprintf("delete:%d", post.ID)),
		),
	)

	msg := tu.Message(tu.ID(int64(config.Conf.TelegramOwner)), post.String())
	msg.ReplyMarkup = keyboard
	bot.SendMessage(context.Background(), msg)

	return nil
}

func DownloadImage(bot *telego.Bot, fileId string) (*string, error) {
	file, err := bot.GetFile(context.Background(), &telego.GetFileParams{FileID: fileId})
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	extension := strings.Split(file.FilePath, ".")[1]
	filename := fmt.Sprintf("./data/images/%s.%s", file.FileUniqueID, extension)

	if !fileExists(filename) {
		downloadURL := bot.FileDownloadURL(file.FilePath)
		err := downloadFile(downloadURL, filename)
		if err != nil {
			return nil, err
		}
	}

	return &filename, nil
}

func IsImageExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return true
	default:
		return false
	}
}

func GetIntArgument(text string) (*int, *merrors.MError) {
	args := strings.Fields(text)[1:]
	if len(args) < 1 {
		return nil, merrors.New(merrors.TelegramArgumentNotProvidedErrorCode, "Missing int argument")
	}

	num, err := strconv.Atoi(args[0])
	if err != nil {
		logger.Error("Delete Post command", err)
		return nil, merrors.New(merrors.CannotConvertToIntErrorCode, "Provided argument is not an int")
	}

	return &num, nil
}
