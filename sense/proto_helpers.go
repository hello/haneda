package sense

import (
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"github.com/hello/haneda/api"
	"github.com/hello/haneda/haneda"
	"github.com/hello/haneda/messeji"
	"time"
)

func GenMesseji(messageId uint64) (*MessageParts, error) {
	header := &haneda.Preamble{}
	header.Type = haneda.Preamble_MESSEJI.Enum()
	header.Id = proto.Uint64(messageId)

	playAudio := &messeji.PlayAudio{}
	playAudio.DurationSeconds = proto.Uint32(10)
	playAudio.FadeInDurationSeconds = proto.Uint32(2)
	playAudio.FadeOutDurationSeconds = proto.Uint32(2)
	playAudio.VolumePercent = proto.Uint32(80)
	playAudio.FilePath = proto.String("/SLPTONES/ST011.RAW")

	msg := &messeji.Message{}
	msg.SenderId = proto.String("go-bench")
	msg.Order = proto.Int64(time.Now().UnixNano())
	msg.Type = messeji.Message_PLAY_AUDIO.Enum()
	msg.PlayAudio = playAudio

	body, _ := proto.Marshal(msg)
	n := SenseId("bench-client")
	mp := &MessageParts{
		Header:  header,
		Body:    body,
		SenseId: n,
	}

	return mp, nil
}

func GenSyncResp(messageId uint64) (*MessageParts, error) {
	syncHeader := &haneda.Preamble{}
	syncHeader.Type = haneda.Preamble_SYNC_RESPONSE.Enum()
	syncHeader.Id = proto.Uint64(messageId)

	syncResponse := &api.SyncResponse{}
	syncResponse.RoomConditions = api.SyncResponse_WARNING.Enum()
	syncResponse.RingTimeAck = proto.String("hi chris")

	body, _ := proto.Marshal(syncResponse)
	n := SenseId("bench-client")
	mp := &MessageParts{
		Header:  syncHeader,
		Body:    body,
		SenseId: n,
	}

	return mp, nil
}

func GenLogs(messageId uint64) (*MessageParts, error) {
	pb := &haneda.Preamble{}
	pb.Type = haneda.Preamble_SENSE_LOG.Enum()
	pb.Id = proto.Uint64(messageId)

	sLog := &api.SenseLog{}
	combined := fmt.Sprintf("Log #%d", messageId)
	sLog.Text = &combined
	n := SenseId("bench-client")
	sLog.DeviceId = proto.String(string(n))
	body, err := proto.Marshal(sLog)
	if err != nil {
		return nil, err
	}
	mp := &MessageParts{
		Header:  pb,
		Body:    body,
		SenseId: n,
	}
	return mp, nil
}

type FileManifestGenerator struct {
	s3service string
}

func (f *FileManifestGenerator) Generate(messageId uint64) (*MessageParts, error) {
	// FileSync.FileManifest.File
	header := &haneda.Preamble{}
	header.Type = haneda.Preamble_FILE_MANIFEST.Enum()
	header.Id = proto.Uint64(messageId)

	fileManifest := &api.FileManifest{}

	fileInfo := fileManifest.GetFileInfo()
	for i := 0; i < 10; i++ {
		fileDownload := &api.FileManifest_FileDownload{}
		fileDownload.Host = proto.String("s3.amazonaws.com")
		fileDownload.Url = proto.String("hello-audio/sleep-tones-raw/ST001.raw")
		file := &api.FileManifest_File{}
		file.DeleteFile = proto.Bool(false)
		file.UpdateFile = proto.Bool(false)
		file.DownloadInfo = fileDownload

		fileInfo = append(fileInfo, file)
	}

	fileManifest.FileInfo = fileInfo

	body, pbErr := proto.Marshal(fileManifest)
	if pbErr != nil {
		return nil, pbErr
	}
	n := SenseId("bench-client")
	mp := &MessageParts{
		Header:  header,
		Body:    body,
		SenseId: n,
	}

	return mp, nil
}

func GenPeriodic(messageId uint64) (*MessageParts, error) {
	header := &haneda.Preamble{}
	header.Type = haneda.Preamble_BATCHED_PERIODIC_DATA.Enum()
	header.Id = proto.Uint64(messageId)

	batched := &api.BatchedPeriodicData{}
	periodic := &api.PeriodicData{}
	periodic.Temperature = proto.Int32(3500)

	n := string("bench-client")
	batched.DeviceId = &n
	batched.FirmwareVersion = proto.Int32(888)
	batched.Data = append(batched.Data, periodic)

	body, pbErr := proto.Marshal(batched)
	if pbErr != nil {
		return nil, pbErr
	}

	mp := &MessageParts{
		Header:  header,
		Body:    body,
		SenseId: SenseId(n),
	}
	return mp, nil
}
