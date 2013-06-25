package petfinder

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	petp "perfectpet4me/pet"
	"strconv"
	"strings"

	"appengine"
	"appengine/urlfetch"
	simplejson "github.com/bitly/go-simplejson"
)

/**********************************************************
                       Constants
*********************************************************/

const API_URL string = "http://api.petfinder.com/"
const API_KEY string = "4b7c674b10cfc0793e935e8abd7346c2"
const SECRET string = "15e0e406bc3bcc284773a57957492133"

/**********************************************************
                   Data Structures
*********************************************************/

type PetFinder struct {
	Token      string
	Expires    string
	w          http.ResponseWriter
	r          *http.Request
	LastOffset int
}

type TokenFetcher struct {
	Petfinder struct {
		Auth struct {
			Token         map[string]string
			Expires       map[string]string
			ExpiresString map[string]string
		}
	}
}

type Pet struct {
	Age         string
	Animal      string
    Breed       []string
    Contact     struct {
		Address1    string
		Address2    string
		City        string
		Email       string
		Fax         string
		State       string
		Zip         string
	}
	Description string
	Id          string
	LastUpdate  string
	Name        string
	Options     []string
	Photos      [](map[string]string)
	Sex         string
	ShelterId   string
	ShelterPetId string
	Size        string
}


/**********************************************************
                       Methods
*********************************************************/
func md5hash(text string) string {
	h := md5.New()
	io.WriteString(h, text)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (pf *PetFinder) SetToken(w http.ResponseWriter, r *http.Request) string {
	method := "auth.getToken"
	errstr := ""
	args := map[string]string{
		"key": API_KEY,
	}
	getstr := pf.RequestBuilder(method, args)

	body, _ := pf.Execute(getstr)

	tf := new(TokenFetcher)
	err := json.Unmarshal(body, &tf)
	if err != nil {
		fmt.Fprintf(w, "%v\n", err.Error())
		errstr = "couldn't unmarshal"
		return errstr
	}

	pf.Token = tf.Petfinder.Auth.Token["$t"]
	pf.Expires = tf.Petfinder.Auth.Expires["$t"]

	return errstr
}

func NewPetFinder(w http.ResponseWriter, r *http.Request) *PetFinder {
	pf := new(PetFinder)
	pf.w = w
	pf.r = r
	err := pf.SetToken(w, r)

	if err != "" {
		return nil
	}

	return pf
}

func (pf *PetFinder) GetPet(animal string, location string) *petp.Pet {
	errstr := ""
	method := "pet.find"
	args := map[string]string{
		"key":      API_KEY,
		"token":    pf.Token,
		"animal":   animal,
		"location": location,
		"count":    "1",
	}
	getstr := pf.RequestBuilder(method, args)

	body, errstr := pf.Execute(getstr)
	if errstr != "" {
		return nil
	}

	j, _ := simplejson.NewJson(body)
	p := new(petp.Pet)
	pjson := j.Get("petfinder").Get("pets").Get("pet")
	p.Id, _ = pjson.Get("id").Get("$t").String()
	p.Name, _ = pjson.Get("name").Get("$t").String()
	p.Age, _ = pjson.Get("age").Get("$t").String()
	p.AnimalType, _ = pjson.Get("animal").Get("$t").String()
	p.Breed, _ = pjson.Get("breeds").Get("breed").Get("$t").String()
	p.Sex, _ = pjson.Get("sex").Get("$t").String()
	p.Description, _ = pjson.Get("description").Get("$t").String()
	x, _ := pjson.Get("media").Get("photos").Get("photo").GetIndex(0).Get("$t").String()
	t, _ := pjson.Get("media").Get("photos").Get("photo").GetIndex(4).Get("$t").String()
	p.PictureURLs[0] = map[string]string{
		"x": x,
		"t": t,
	}

	return p
}

func (pf *PetFinder) GetPets(animal string, location string, numResults int) (pets []*Pet) {
    c := appengine.NewContext(pf.r)

	method := "pet.find"
	args := map[string]string{
		"key":      API_KEY,
		"token":    pf.Token,
		"animal":   animal,
		"location": location,
		"count":    strconv.Itoa(numResults),
	}
	getstr := pf.RequestBuilder(method, args)

	body, errstr := pf.Execute(getstr)
	if errstr != "" {
		return nil
	}

	j, _ := simplejson.NewJson(body)
    var err error

	for i := 0; i < numResults; i++ {
		p := new(Pet)
		pjson := j.Get("petfinder").Get("pets").Get("pet").GetIndex(i)

        // age
        if age, ok := pjson.CheckGet("age"); ok {
            if p.Age, err = age.Get("$t").String(); err != nil {
                c.Errorf("age: %v\n", err)
            }
        }
        // animal
        if animal, ok := pjson.CheckGet("animal"); ok {
            if p.Animal, err = animal.Get("$t").String(); err != nil {
                c.Errorf("animal: %v\n", err)
            }
        }
        // breeds
        if breeds, ok := pjson.CheckGet("breeds"); ok {
            if breed, ok := breeds.CheckGet("breed"); ok {
                barray, _ := breed.Array()
                if barray == nil {
                    if breedstr, ok := breed.CheckGet("$t"); ok {
                        if temp, err := breedstr.String(); err != nil {
                            c.Errorf("one breed: %v\n", err)
                        } else {
                            p.Breed = append(p.Breed, temp)
                        }
                    }
                } else {
                    for i := range barray {
                        if temp, err := breed.GetIndex(i).Get("$t").String(); err != nil {
                            c.Errorf("multiple breeds: %v\n", err)
                        } else {
                            p.Breed = append(p.Breed, temp)
                        }
                    }
                }
            } else {
                c.Errorf("breed: %v\n", err)
            }
        }
        // contact
        if contact, ok := pjson.CheckGet("contact"); ok {
            // address1
            if address1, ok := contact.CheckGet("address1"); ok {
                if p.Contact.Address1, err = address1.Get("$t").String(); err != nil {
                    c.Errorf("address1: %v\n", err)
                }
            }
            // address2
            if address2, ok := contact.CheckGet("address2"); ok {
                if p.Contact.Address2, err = address2.Get("$t").String(); err != nil {
                    c.Errorf("address2: %v\n", err)
                }
            }
            // city
            if city, ok := contact.CheckGet("city"); ok {
                if p.Contact.City, err = city.Get("$t").String(); err != nil {
                    c.Errorf("city: %v\n", err)
                }
            }
            // email
            if email, ok := contact.CheckGet("email"); ok {
                if p.Contact.Email, err = email.Get("$t").String(); err != nil {
                    c.Errorf("email: %v\n", err)
                }
            }
            // fax
            if fax, ok := contact.CheckGet("fax"); ok {
                if p.Contact.Fax, err = fax.Get("$t").String(); err != nil {
                    c.Errorf("fax: %v\n", err)
                }
            }
            // state
            if state, ok := contact.CheckGet("state"); ok {
                if p.Contact.State, err = state.Get("$t").String(); err != nil {
                    c.Errorf("state: %v\n", err)
                }
            }
            // zip
            if zip, ok := contact.CheckGet("zip"); ok {
                if p.Contact.Zip, err = zip.Get("$t").String(); err != nil {
                    temp, err := zip.Get("$t").Float64()
                    p.Contact.Zip = strconv.Itoa(int(temp))
                    if err != nil {
                        c.Errorf("zip: %v\n", err)
                    }
                }
            }
        }
        // description
        if description, ok := pjson.CheckGet("description"); ok {
            if p.Description, err = description.Get("$t").String(); err != nil {
                c.Errorf("description: %v\n", err)
            }
        }
        // id
        if id, ok := pjson.CheckGet("id"); ok {
            if p.Id, err = id.Get("$t").String(); err != nil {
                temp, err := id.Get("$t").Float64()
                p.Id = strconv.Itoa(int(temp))
                if err != nil {
                    c.Errorf("id: %v\n", err)
                }
            }
        }
        // lastupdate
        if lastupdate, ok := pjson.CheckGet("lastupdate"); ok {
            if p.LastUpdate, err = lastupdate.Get("$t").String(); err != nil {
                c.Errorf("lastupdate: %v\n", err)
            }
        }
        // name
        if name, ok := pjson.CheckGet("name"); ok {
            if p.Name, err = name.Get("$t").String(); err != nil {
                c.Errorf("name: %v\n", err)
            }
        }
        // options
        if options, ok := pjson.CheckGet("options"); ok {
            if option, ok := options.CheckGet("option"); ok {
                oarray, _ := option.Array()
                if oarray == nil {
                    if optionstr, ok := option.CheckGet("$t"); ok {
                        if temp, err := optionstr.String(); err != nil {
                            c.Errorf("one option: %v\n", err)
                        } else {
                            p.Options = append(p.Options, temp)
                        }
                    }
                } else {
                    for i := range oarray {
                        if temp, err := option.GetIndex(i).Get("$t").String(); err != nil {
                            c.Errorf("multiple options: %v\n", err)
                        } else {
                            p.Options = append(p.Options, temp)
                        }
                    }
                }
            } else {
                c.Errorf("option: %v\n", err)
            }
        }
        // photos
        if photos, ok := pjson.Get("media").CheckGet("photos"); ok {
            if photo, ok := photos.CheckGet("photo"); ok {
                parray, _ := photo.Array()
                if parray == nil {
                    c.Errorf("empty photo array!!!: %v\n", err)
                } else {
                    for i := range parray {
                        var size, url string
                        var id int
                        if id, err = photo.GetIndex(i).Get("@id").Int(); err != nil {
                            c.Errorf("photo id: %v\n", err)
                        }
                        if size, err = photo.GetIndex(i).Get("@size").String(); err != nil {
                            c.Errorf("photo size: %v\n", err)
                        }
                        if url, err = photo.GetIndex(i).Get("$t").String(); err != nil {
                            c.Errorf("photo url: %v\n", err)
                        }
                        if err == nil {
                            p.Photos = append(p.Photos, map[string]string{
                                "id": strconv.Itoa(id),
                                "size": size,
                                "url": url,
                                } )
                        }
                    }
                }
            } else {
                c.Errorf("photo: %v\n", err)
            }
        }
        // sex
        if sex, ok := pjson.CheckGet("sex"); ok {
            if p.Sex, err = sex.Get("$t").String(); err != nil {
                c.Errorf("sex: %v\n", err)
            }
        }
        // shelterid
        if shelterid, ok := pjson.CheckGet("shelterid"); ok {
            if p.ShelterId, err = shelterid.Get("$t").String(); err != nil {
                c.Errorf("shelterid: %v\n", err)
            }
        }
        // shelterpetid
        if shelterpetid, ok := pjson.CheckGet("shelterpetid"); ok {
            if p.ShelterPetId, err = shelterpetid.Get("$t").String(); err != nil {
                c.Errorf("shelterpid: %v\n", err)
            }
        }
        // size
        if size, ok := pjson.CheckGet("size"); ok {
            if p.Size, err = size.Get("$t").String(); err != nil {
                c.Errorf("size: %v\n", err)
            }
        }

        pets = append(pets, p)
	}
	return
}

func (pf *PetFinder) RequestBuilder(apicall string, args map[string]string) string {
	split := strings.Split(apicall, ".")
	getstr := ""
	switch objtype := split[0]; objtype {
	case "pet":
		switch method := split[1]; method {
		case "get":
		case "getRandom":
		case "find":
			getstr = "key=" + args["key"] + "&token=" + args["token"] + "&animal=" + args["animal"] + "&location=" + args["location"] + "&count=" + args["count"] + "&format=json"
			signature := md5hash(SECRET + getstr)
			getstr = apicall + "?" + getstr + "&sig=" + signature
		}
	case "breed":
		switch method := split[1]; method {
		case "list":
		}
	case "shelter":
		switch method := split[1]; method {
		case "find":
		case "get":
		case "getPets":
		case "listByBreed":
		}
	case "auth":
		switch method := split[1]; method {
		case "getToken":
			getstr = "key=" + args["key"] + "&format=json"
			signature := md5hash(SECRET + getstr)
			getstr = apicall + "?" + getstr + "&sig=" + signature
		}
	}

	return API_URL + getstr
}

func (pf *PetFinder) Execute(getstr string) ([]byte, string) {
	errstr := ""

	c := appengine.NewContext(pf.r)
	client := urlfetch.Client(c)

	resp, err := client.Get(getstr)
	if err != nil {
		errstr = "couldn't open client"
		return nil, errstr
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		errstr = "couldn't read body"
		return body, errstr
	}

	return body, errstr
}
