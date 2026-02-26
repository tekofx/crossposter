package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/tekofx/crossposter/internal/config"
	"github.com/tekofx/crossposter/internal/logger"
	"github.com/tekofx/crossposter/internal/model"

	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func SendMessageToOwner(ctx *th.Context, text string) (*telego.Message, error) {
	msg, err := ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(int64(config.Conf.TelegramOwner)),
		text,
	).WithParseMode(telego.ModeMarkdownV2))

	if err != nil {
		return nil, err
	}
	return msg, nil
}

func SendMessageToOwnerUsingBot(bot *telego.Bot, text string) (*telego.Message, error) {
	msg, err := bot.SendMessage(context.Background(), tu.Message(
		tu.ID(int64(config.Conf.TelegramOwner)),
		text,
	).WithParseMode(telego.ModeMarkdownV2))

	if err != nil {
		return nil, err
	}
	return msg, nil
}

func SendMediaGroupByFileIDs(bot *telego.Bot, chatID int64, post *model.Post) error {
	var media []telego.InputMedia

	for i, image := range post.Images {

		downloadedImage, err := os.Open(image.Filename)
		if err != nil {
			logger.Error("Error opening file:", err)
			return err
		}
		defer downloadedImage.Close()

		inputFile := telego.InputFile{
			File: downloadedImage,
		}

		photo := &telego.InputMediaPhoto{
			Type:  "photo",
			Media: inputFile,
		}
		if i == 0 {
			photo.Caption = post.Text
		}
		media = append(media, photo)
	}

	params := telego.SendMediaGroupParams{
		ChatID: telego.ChatID{ID: chatID},
		Media:  media,
	}

	_, err := bot.SendMediaGroup(context.Background(), &params)
	return err
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
