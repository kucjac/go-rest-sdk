package restsdk

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestStatusStringer(t *testing.T) {
	var status ResponseStatus
	Convey("While having the status of value '1'", t, func() {
		status = ResponseStatus(1)
		Convey("The String method should return 'ok'", func() {
			So(status.String(), ShouldEqual, "ok")
		})
	})
	Convey("While having the status of value '2'", t, func() {
		status = ResponseStatus(2)

		Convey("The string method should return 'error'", func() {
			So(status.String(), ShouldEqual, "error")
		})
	})
	Convey("For any other value of the response status", t, func() {
		status = ResponseStatus(5)
		Convey("The String() method should return 'unknown'", func() {
			So(status.String(), ShouldEqual, "unknown")
		})
	})

}

func TestMarshalStatus(t *testing.T) {
	var status ResponseStatus

	Convey("Having the ResponseStatus 'ok' - of value 1", t, func() {
		status = ResponseStatus(1)

		Convey("It should Marshal as a string value equal to \"ok\"", func() {
			marshaled, err := json.Marshal(&status)
			So(err, ShouldBeNil)
			So(string(marshaled), ShouldEqual, "\"ok\"")
		})
	})
}

func TestUnmarshalStatus(t *testing.T) {
	type wrapper struct {
		Status ResponseStatus `json:"status"`
	}

	Convey("Having the JSON {\"status\" : \"ok\"}", t, func() {
		var statusOk string = `{"status" :"ok"}`

		Convey("And unmarshaling it into 'ResponseStatus", func() {
			var okWrapper wrapper
			err := json.Unmarshal([]byte(statusOk), &okWrapper)

			So(err, ShouldBeNil)

			Convey("The response status value should be equal 1", func() {
				So(okWrapper.Status, ShouldEqual, ResponseStatus(1))
			})
		})
	})

	Convey("Having the JSON {\"status\" : \"error\"}", t, func() {
		var statusError string = `{"status": "error"}`

		Convey("And unmarshaling it into 'ResponseStatus", func() {
			var errorWrapper wrapper
			err := json.Unmarshal([]byte(statusError), &errorWrapper)

			So(err, ShouldBeNil)

			Convey("The response status value should be equal 2", func() {
				So(errorWrapper.Status, ShouldEqual, ResponseStatus(2))
			})
		})
	})

	Convey("Having the JSON {\"status\" : \"unknown\"}", t, func() {
		var statusUnknown string = `{"status": "unknown"}`

		Convey("And unmarshaling it into 'ResponseStatus", func() {
			var unknownWrapper wrapper

			err := json.Unmarshal([]byte(statusUnknown), &unknownWrapper)
			So(err, ShouldBeNil)

			Convey("The response status value should be equal 0", func() {
				So(unknownWrapper.Status, ShouldEqual, ResponseStatus(0))
			})
		})
	})

	Convey("Having the JSON {\"status\": 1.2}", t, func() {
		var statusIncorrect string = `{"status": 1.2}`

		Convey("Unmarshaling it into 'ResponseStatus'", func() {
			var incorrectWrapper wrapper

			err := json.Unmarshal([]byte(statusIncorrect), &incorrectWrapper)

			Convey("It should throws an error", func() {
				So(err, ShouldBeError)
			})
		})
	})
}
