package nativeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/navidrome/navidrome/log"
	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/utils/req"
)

type BeetsItem struct {
	Results []struct {
		Artpath              string  `json:"artpath"`
		Year                 int     `json:"year"`
		Albumtype            string  `json:"albumtype"`
		R128TrackGain        any     `json:"r128_track_gain"`
		Day                  int     `json:"day"`
		Disc                 int     `json:"disc"`
		Albumartist          string  `json:"albumartist"`
		MbReleasegroupid     string  `json:"mb_releasegroupid"`
		InitialKey           any     `json:"initial_key"`
		AcoustidID           string  `json:"acoustid_id"`
		MbAlbumid            string  `json:"mb_albumid"`
		Month                int     `json:"month"`
		Track                int     `json:"track"`
		Style                string  `json:"style"`
		RgAlbumGain          any     `json:"rg_album_gain"`
		Album                string  `json:"album"`
		DiscogsArtistid      int     `json:"discogs_artistid"`
		Tracktotal           int     `json:"tracktotal"`
		MbAlbumartistid      string  `json:"mb_albumartistid"`
		Country              string  `json:"country"`
		Channels             int     `json:"channels"`
		Arranger             string  `json:"arranger"`
		Comp                 int     `json:"comp"`
		Bitrate              int     `json:"bitrate"`
		Length               float64 `json:"length"`
		Lyricist             string  `json:"lyricist"`
		RgAlbumPeak          any     `json:"rg_album_peak"`
		MbArtistid           string  `json:"mb_artistid"`
		Trackdisambig        string  `json:"trackdisambig"`
		OriginalMonth        int     `json:"original_month"`
		AcoustidFingerprint  string  `json:"acoustid_fingerprint"`
		OriginalDay          int     `json:"original_day"`
		ComposerSort         string  `json:"composer_sort"`
		AlbumartistCredit    string  `json:"albumartist_credit"`
		MbWorkid             string  `json:"mb_workid"`
		MbReleasetrackid     string  `json:"mb_releasetrackid"`
		Mtime                float64 `json:"mtime"`
		Work                 string  `json:"work"`
		AlbumartistSort      string  `json:"albumartist_sort"`
		DiscogsLabelid       int     `json:"discogs_labelid"`
		Bpm                  int     `json:"bpm"`
		Language             string  `json:"language"`
		DataSource           string  `json:"data_source"`
		ArtistCredit         string  `json:"artist_credit"`
		Format               string  `json:"format"`
		Composer             string  `json:"composer"`
		Disctotal            int     `json:"disctotal"`
		Title                string  `json:"title"`
		Grouping             string  `json:"grouping"`
		Added                float64 `json:"added"`
		Media                string  `json:"media"`
		Artist               string  `json:"artist"`
		Albumdisambig        string  `json:"albumdisambig"`
		Isrc                 string  `json:"isrc"`
		DiscogsAlbumid       int     `json:"discogs_albumid"`
		Disctitle            string  `json:"disctitle"`
		Lyrics               string  `json:"lyrics"`
		Albumstatus          string  `json:"albumstatus"`
		Albumtypes           any  `json:"albumtypes"`
		Releasegroupdisambig string  `json:"releasegroupdisambig"`
		Comments             string  `json:"comments"`
		Encoder              string  `json:"encoder"`
		Catalognum           string  `json:"catalognum"`
		R128AlbumGain        any     `json:"r128_album_gain"`
		Label                string  `json:"label"`
		Bitdepth             int     `json:"bitdepth"`
		RgTrackGain          any     `json:"rg_track_gain"`
		ID                   int     `json:"id"`
		Script               string  `json:"script"`
		ArtistSort           string  `json:"artist_sort"`
		TrackAlt             string  `json:"track_alt"`
		Genre                string  `json:"genre"`
		OriginalYear         int     `json:"original_year"`
		WorkDisambig         string  `json:"work_disambig"`
		AlbumID              int     `json:"album_id"`
		MbTrackid            string  `json:"mb_trackid"`
		Asin                 string  `json:"asin"`
		RgTrackPeak          any     `json:"rg_track_peak"`
		Samplerate           int     `json:"samplerate"`
		Size                 int     `json:"size"`
	} `json:"results"`
}

func deleteSong(ds model.DataStore) http.HandlerFunc {
	type deleteSongPayload struct {
		Ids []string `json:"ids"`
		User string `json:"user"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// todo: tests, use proper url parsing lib
		var payload deleteSongPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		for _, id := range payload.Ids {
			println(id)
			mf, err := ds.MediaFile(ctx).Get(id)
			if err != nil {
				log.Error(err)
			}
			println(mf.Artist)
			println(mf.Title)
			// todo set this base from env variable
			//baseUrl := "http://127.0.0.1:8337"
			// get docker port
		    baseUrl := fmt.Sprintf("http://beets%s:8337", payload.User)
			queryEndPoint := "/item/query/"
			queryStr := fmt.Sprintf("artist:%s/title:%s/album:%s/user:%s", mf.Artist, mf.Title, mf.Album, payload.User)
			url := baseUrl + queryEndPoint + queryStr
			fmt.Printf("query url: %s\n", url)
			resp, err := http.Get(url) // nolint
			if err != nil {
				log.Error(err)
			    http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Error(err)
			    http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			sb := string(body)
			var beetsItem BeetsItem
			err = json.Unmarshal([]byte(sb), &beetsItem)
			if err != nil {
				log.Error(err)
			    http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			length := len(beetsItem.Results)
			log.Info("following query string matched ", length, "entries: ", queryStr)
			for _, item := range beetsItem.Results {
				log.Info("deleting: ", item.Artist, " : ", item.Title, " id: ", item.ID)

				deleteStr := fmt.Sprintf("/item/%d", item.ID)

				deleteFromDiskStr := "?delete"

				del_url := baseUrl + deleteStr + deleteFromDiskStr
				// Create request
				req, err := http.NewRequest(http.MethodDelete, del_url, nil)
				log.Info("delete request: ", req)
				if err != nil {
					log.Error(err)
				    http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// Fetch Request
				client := &http.Client{}
				del_resp, del_err := client.Do(req)
				if del_err != nil {
					log.Error(del_err)
				    http.Error(w, del_err.Error(), http.StatusInternalServerError)
					return
				}
				defer del_resp.Body.Close()

				// read the response body
				del_body, err := io.ReadAll(del_resp.Body)
				if err != nil {
					log.Error(err)
				    http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// print the response body
				fmt.Println(string(del_body))
				log.Info("del body: ", del_body)
			}
		}
	}
}

func getBeetTrack(ds model.DataStore) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// todo: tests, use proper url parsing lib
	    p := req.Params(r)
	    id, _ := p.String("id")
		user, _ := p.String("user")

		log.Info("got parameters", "id", id, "user", user)
		ctx := r.Context()
		mf, err := ds.MediaFile(ctx).Get(id)
		if err != nil {
			log.Error(err)
		}
		if mf == nil {
			log.Error("no mediafile for id", "id", id)
		}
		log.Info("mediafile", "mf", mf)
		// todo set this base from env variable
		//baseUrl := "http://127.0.0.1:8337"
		baseUrl := fmt.Sprintf("http://beets:8337")
		queryEndPoint := "/item/query/"
		queryStr := fmt.Sprintf("artist:%s/title:%s/user:%s", mf.Artist, mf.Title, user)
		url := baseUrl + queryEndPoint + queryStr
		fmt.Printf("query url: %s\n", url)
		resp, err := http.Get(url) // nolint
		if err != nil {
			log.Error(err)
		    http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
		    http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sb := string(body)
		var beetsItem BeetsItem
		err = json.Unmarshal([]byte(sb), &beetsItem)
		if err != nil {
			log.Error(err)
		    http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(beetsItem)
	}
}
