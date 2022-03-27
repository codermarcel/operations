package src

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testService = `
version: beta

profiles:
  dev:
    name: test-app-dev
    max_cpu: 0.25
  staging:
    name: test-app-staging
    max_cpu: 0.5
  prod:
    name: test-app-prod
    max_cpu: 1
`

var testTemplate1 = `
apiVersion: v1
kind: Service
metadata:
	name: {{ .name }}
spec:
	selector:
	app: {{ .max_cpu }}
	ports:
	- port: 80
		targetPort: go-app
	type: ClusterIP
`

var dev = `
apiVersion: v1
kind: Service
metadata:
	name: test-app-dev
spec:
	selector:
	app: 0.25
	ports:
	- port: 80
		targetPort: go-app
	type: ClusterIP
`
var staging = `
apiVersion: v1
kind: Service
metadata:
	name: test-app-staging
spec:
	selector:
	app: 0.5
	ports:
	- port: 80
		targetPort: go-app
	type: ClusterIP
`

var prod = `
apiVersion: v1
kind: Service
metadata:
	name: test-app-prod
spec:
	selector:
	app: 1
	ports:
	- port: 80
		targetPort: go-app
	type: ClusterIP
`

var testTemplate2 = `
name: {{ envOr "someenv" "default" }}
`

var expectedEnv = `
name: value
`

var expectedEnvDefault = `
name: default
`

func TestProfileBasic(t *testing.T) {
	assert := assert.New(t)

	svc, err := NewService(ServiceFromBytes([]byte(testService)))
	assert.Nil(err)

	tmpl, err := NewBaseTemplate(TemplateFromString(testTemplate1))
	assert.Nil(err)

	assert.Equal(3, len(svc.Profiles))
	assert.True(svc.HasProfile("dev"))
	assert.True(svc.HasProfile("staging"))
	assert.True(svc.HasProfile("prod"))

	resBufDev := &bytes.Buffer{}
	err = tmpl.ExecuteTo(svc.Profiles["dev"], resBufDev)
	assert.Nil(err)
	assert.NotNil(resBufDev)
	assert.Equal(dev, resBufDev.String())

	resBufStaging := &bytes.Buffer{}
	err = tmpl.ExecuteTo(svc.Profiles["staging"], resBufStaging)
	assert.Nil(err)
	assert.NotNil(resBufStaging)
	assert.Equal(staging, resBufStaging.String())

	resBufProd := &bytes.Buffer{}
	err = tmpl.ExecuteTo(svc.Profiles["prod"], resBufProd)
	assert.Nil(err)
	assert.NotNil(resBufProd)
	assert.Equal(prod, resBufProd.String())
}

func TestProfileEnvOr(t *testing.T) {
	assert := assert.New(t)

	svc, err := NewService(ServiceFromBytes([]byte(testService)))
	assert.Nil(err)

	tmpl, err := NewBaseTemplate(TemplateFromString(testTemplate2))
	assert.Nil(err)

	envName := "someenv"

	os.Unsetenv(envName)

	first := &bytes.Buffer{}
	err = tmpl.ExecuteTo(svc.Profiles["dev"], first)
	assert.Nil(err)
	assert.NotNil(first)
	assert.Equal(expectedEnvDefault, first.String())

	os.Setenv(envName, "value")

	second := &bytes.Buffer{}
	err = tmpl.ExecuteTo(svc.Profiles["dev"], second)
	assert.Nil(err)
	assert.NotNil(second)
	assert.Equal(expectedEnv, second.String())

	assert.NotEqual(expectedEnv, expectedEnvDefault)
}
func TestProfileFileOutput(t *testing.T) {
	assert := assert.New(t)

	svc, err := NewService(ServiceFromBytes([]byte(testService)))
	assert.Nil(err)

	tmpl, err := NewBaseTemplate(TemplateFromString(testTemplate2))
	assert.Nil(err)

	name := "_testfile_profile.yaml"

	f, err := os.Create(name)
	assert.Nil(err)
	defer os.Remove(name)

	err = tmpl.ExecuteTo(svc.Profiles["dev"], f)
	assert.Nil(err)
	assert.NotNil(f)
}
