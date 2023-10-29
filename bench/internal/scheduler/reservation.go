package scheduler

import (
	"log"
	"time"

	"github.com/biogo/store/interval"
)

// FIXME: 同時配信枠数は、少なくとも２は想定しておいたほうがいい
// FIXME: 枠数に合わせる
const hotThreshold = 1

func ConvertFromIntInterface(i []interval.IntInterface) ([]*Reservation, error) {
	reservations := make([]*Reservation, len(i))
	for idx, ii := range i {
		reservation, ok := ii.(*Reservation)
		if !ok {
			log.Println("failed to convert reservation")
			continue
		}
		reservations[idx] = reservation
	}

	return reservations, nil
}

type Reservation struct {
	// NOTE: id は、webappで割り振られるIDではなく、ReservationSchedulerが管理する上で利用するもの
	id          int
	Title       string
	Description string
	Tags        []int
	StartAt     int64
	EndAt       int64
}

// FIXME: id, UserNameなど古い引数を廃止
// 初期データ生成スクリプト側修正後実施
func mustNewReservation(id int, UserName string, title string, description string, startAtStr string, endAtStr string) *Reservation {
	startAt, err := time.Parse("2006-01-02 15:04:05", startAtStr)
	if err != nil {
		log.Fatalln(err)
	}
	endAt, err := time.Parse("2006-01-02 15:04:05", endAtStr)
	if err != nil {
		log.Fatalln(err)
	}

	// FIXME: タグの採番がおかしくてwebappでエラーが出る
	// tagCount := rand.Intn(10)
	reservation := &Reservation{
		id:          id,
		Title:       title,
		Description: description,
		StartAt:     startAt.Unix(),
		EndAt:       endAt.Unix(),
	}
	// for i := 0; i < tagCount; i++ {
	// 	tagIdx := rand.Intn(len(tagPool))
	// 	reservation.Tags = append(reservation.Tags, tagIdx)
	// }

	return reservation
}

func (r *Reservation) Overlap(interval interval.IntRange) bool {
	if interval.Start == interval.End {
		// 区間の開始と終了が同じである場合、予約の中に含まれるならオーバーラップと判定させる
		return r.StartAt <= int64(interval.Start) && r.EndAt >= int64(interval.Start)
	}
	// NOTE: 指定区間の外側についてexclusiveな判定を行う
	//       指定区間の内側についてinclusiveな判定を行う
	if r.StartAt >= int64(interval.End) {
		// 予約開始が指定区間の終了以上である場合は含めない
		return false
	}
	if r.EndAt <= int64(interval.Start) {
		// 予約終了が指定区間の開始以下である場合は含めない
		return false
	}
	return r.EndAt >= int64(interval.Start) && r.StartAt <= int64(interval.End)
}
func (r *Reservation) ID() uintptr { return uintptr(r.id) }
func (r *Reservation) Range() interval.IntRange {
	return interval.IntRange{Start: int(r.StartAt), End: int(r.EndAt)}
}
