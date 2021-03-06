package main

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Layer 자료구조
type Layer struct {
	Name   string `json:"name" bson:"name"`     // Layer 이름
	Color  string `json:"color" bson:"color"`   // Layer 컬러
	Order  int    `json:"order" bson:"order"`   // Layer 순서
	Hidden bool   `json:"hidden" bson:"hidden"` // Layer 숨김 속성
}

// Schedule 자료구조
type Schedule struct {
	ID         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Collection string        `json:"collection" bson:"collection"` // 사용자,장비명,회의실 이름이 될 수 있다
	Title      string        `json:"title" bson:"title"`           // 스케쥴의 title
	Start      string        `json:"start" bson:"start"`           // 스케쥴 시작 시간
	Startnum   int64         `json:"startnum" bson:"startnum"`     // 스케쥴 시작 시간을 Excel 5digit general number으로 변환한 값. 이하 '엑셀날짜'로 표기한다.
	End        string        `json:"end" bson:"end"`               // 스케쥴 끝나는 시간
	Endnum     int64         `json:"endnum" bson:"endnum"`         // 스케쥴 끝나는 시간 Excel 5digit general number으로 변환한 값
	Color      string        `json:"color" bson:"color"`           //#FF3366, #ff3366, #f36, #F36
	Layer      string        `json:"layer" bson:"layer"`           // 스케쥴이 속한 레이어의 이름
}

// CheckError 매소드는 Schedule 자료구조에 에러가 있는지 체크한다.
func (s Schedule) CheckError() error {
	if s.Collection == "" {
		return errors.New("Collection 이 빈 문자열 입니다")
	}
	if s.Title == "" {
		return errors.New("Title 이 빈 문자열 입니다")
	}
	if s.Layer == "" {
		return errors.New("Layer 이름이 빈 문자열 입니다")
	}

	if s.Start == "" {
		return errors.New("Start 시간이 빈 문자열 입니다")
	}
	if s.End == "" {
		return errors.New("End 시간이 빈 문자열 입니다")
	}
	startTime, err := time.Parse(time.RFC3339, s.Start)
	if err != nil {
		return err
	}
	endTime, err := time.Parse(time.RFC3339, s.End)
	if err != nil {
		return err
	}
	// end가 start 시간보다 큰지 체크하는 부분
	if !endTime.After(startTime) {
		return errors.New("끝시간이 시작시간보다 작습니다")
	}
	if s.Color != "" {
		if !regexWebColor.MatchString(s.Color) {
			return errors.New("#FF0011 형식의 문자열이 아닙니다")
		}
	}
	return nil
}

// ToUTC 메소드는 스케줄을 받아서 내부 시간을 UTC로 바꾼다.
func (s *Schedule) ToUTC() error {
	startTime, err := time.Parse(time.RFC3339, s.Start)
	if err != nil {
		return err
	}
	endTime, err := time.Parse(time.RFC3339, s.End)
	if err != nil {
		return err
	}
	s.Start = startTime.UTC().Format(time.RFC3339)
	s.End = endTime.UTC().Format(time.RFC3339)
	return nil
}

// SetTimeNum 메소드는 스케줄에 엑셀시간을 셋팅한다.
func (s *Schedule) SetTimeNum() error {
	startNum, err := TimeToNum(s.Start)
	if err != nil {
		return err
	}
	s.Startnum = startNum
	endNum, err := TimeToNum(s.End)
	if err != nil {
		return err
	}
	s.Endnum = endNum
	return nil
}

// CheckError 매소드는 Layer 자료구조에 에러가 있는지 체크한다.
func (l Layer) CheckError() error {
	if l.Name == "" {
		return errors.New("Name 이 빈 문자열 입니다")
	}
	if l.Color != "" {
		if !regexWebColor.MatchString(l.Color) {
			return errors.New("#FF0011 형식의 컬러가 아닙니다")
		}
	}
	return nil
}
