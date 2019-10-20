package binghomeimage

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
)

func GetHttpClient() *http.Client {
	// Cookie handle
	jar, _ := cookiejar.New(nil)

	return &http.Client{Jar: jar}
}

func GetHttpRequest(hMethod, reqUrl string) (*http.Request, error) {
	httpRequest, err := http.NewRequest(strings.ToUpper(hMethod), reqUrl, nil)
	if err != nil {
		return nil, err
	}
	return httpRequest, nil
}

func DownLoadBingHomeImage(bingUrl string) {
	hNode, err := httpGet(bingUrl)
	if err != nil {
		_ = fmt.Errorf("%s", err)
	}
	headNode, _, err := findHeadBody(hNode)
	attrMap := map[string]string{"id": "bgLink", "rel": "preload", "href": "", "as": "image", "iCount": "0"}
	_, err = FindTargetNode(headNode, attrMap)
	if err != nil {
		log.Fatal(err)
	}
	dUrl, err := getDownloadUrl(attrMap["href"], bingUrl)
	fmt.Printf("Get Image Url: <%s> success\n", dUrl)
	if err != nil {
		log.Fatal(err)
	}
	fileName := strings.Split(dUrl, "=")[1]
	err = downFile(dUrl, fileName)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Download <%s> success!\nFile Name: <%s>", dUrl, fileName)
	}
}

func httpGet(url string) (*html.Node, error) {
	client := GetHttpClient()
	req, err := GetHttpRequest("GET", url)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	} else {
		doc, err := html.Parse(resp.Body)
		if err != nil {
			return nil, err
		}
		return doc, nil
	}
}

func findHeadBody(hNode *html.Node) (*html.Node, *html.Node, error) {
	if hNode.Type == html.DocumentNode {
		hNode = hNode.FirstChild.NextSibling
	} else {
		return nil, nil, fmt.Errorf("can find html Document")
	}
	headNode := hNode.FirstChild
	if headNode.Type == html.ElementNode && headNode.Data == "head" {
		bodyNode := headNode.NextSibling
		if bodyNode.Type == html.ElementNode && bodyNode.Data == "body" {
			return headNode, bodyNode, nil
		} else {
			return nil, nil, fmt.Errorf("can find html body")
		}
	} else {
		return nil, nil, fmt.Errorf("can find html head")
	}
}

func FindTargetNode(node *html.Node, tMap map[string]string) (*html.Node, error) {
	loopNode := node.FirstChild
	for loopNode != nil {
		if loopNode.Data == "link" {
			if len(loopNode.Attr) == 4 {
				for _, v := range loopNode.Attr {
					val, ok := tMap[v.Key]
					if ok {
						if val == v.Val {
							iNum, _ := strconv.Atoi(tMap["iCount"])
							tMap["iCount"] = strconv.Itoa(iNum + 1)
						} else if v.Key == "href" {
							tMap[v.Key] = v.Val
						}
					}
				}
				if tMap["iCount"] == "3" && len(tMap["href"]) > 0 {
					return loopNode, nil
				}
			}
		}
		loopNode = loopNode.NextSibling
	}
	return nil, fmt.Errorf("can find target node")
}

func getDownloadUrl(inStr string, baseUrl string) (string, error) {
	strSlice := strings.Split(inStr, "&")
	fragSlice := strings.Split(strSlice[0], "_")
	if len(fragSlice) >= 2 {
		if strings.Contains(fragSlice[len(fragSlice)-1], ".jpg") {
			fragSlice[len(fragSlice)-1] = "UHD.jpg"
			fragSlice[0] = baseUrl + fragSlice[0]
			retStr := strings.Join(fragSlice, "_")
			return retStr, nil
		} else {
			return "", fmt.Errorf(".jpg not in url fragment <%s>", strSlice[0])
		}
	} else {
		return "", fmt.Errorf("url fragment error <%s>", strSlice[0])
	}
}

func downFile(dUrl string, fName string) error {
	// allocate file
	file, err := os.Create(fName)
	if err != nil {
		return err
	}
	defer file.Close()
	client := GetHttpClient()
	req, err := GetHttpRequest("GET", dUrl)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("response status error:%d", resp.StatusCode)
	}
	buf := make([]byte, 8192)
	var writeIndex int64 = 0
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			writeSize, err := file.WriteAt(buf[0:n], writeIndex)
			if err != nil {
				return err
			}
			writeIndex += int64(writeSize)
		}
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
	}
	return nil
}
