package skel

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/deislabs/ratify/pkg/common"
	"github.com/deislabs/ratify/pkg/common/plugin"
	"github.com/deislabs/ratify/pkg/ocispecs"
	"github.com/deislabs/ratify/pkg/referrerstore"
	sp "github.com/deislabs/ratify/pkg/referrerstore/plugin"
	"github.com/deislabs/ratify/pkg/referrerstore/types"
	"github.com/deislabs/ratify/pkg/utils"
	"github.com/opencontainers/go-digest"
)

type pcontext struct {
	GetEnviron func(string) string
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
}

// TODO use pointers to avoid copy
type ListReferrers func(args *CmdArgs, subjectReference common.Reference, artifactTypes []string, nextToken string) (referrerstore.ListReferrersResult, error)
type GetBlobContent func(args *CmdArgs, subjectReference common.Reference, digest digest.Digest) ([]byte, error)
type GetReferenceManifest func(args *CmdArgs, subjectReference common.Reference, digest digest.Digest) (ocispecs.ReferenceManifest, error)

type CmdArgs struct {
	Version    string
	Subject    string
	Args       string
	subjectRef common.Reference
	StdinData  []byte
}

// PluginMain is the core "main" for a plugin which includes automatic error handling.
func PluginMain(name, version string, listReferrers ListReferrers, getBlobContent GetBlobContent, getRefManifest GetReferenceManifest, supportedVersions []string) {
	if e := (&pcontext{
		GetEnviron: os.Getenv,
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
	}).pluginMainCore(name, version, listReferrers, getBlobContent, getRefManifest, supportedVersions); e != nil {
		if err := e.Print(); err != nil {
			log.Print("Error writing error result to stdout: ", err)
		}
		os.Exit(1)
	}
}

func (c *pcontext) pluginMainCore(name, version string, listReferrers ListReferrers, getBlobContent GetBlobContent, getRefManifest GetReferenceManifest, supportedVersions []string) *plugin.Error {
	cmd, cmdArgs, err := c.getCmdArgsFromEnv()
	if err != nil {
		// TODO about string
		return err
	}

	if err = validateVersion(cmdArgs.Version, supportedVersions); err != nil {
		return err
	}

	if err = validateConfig(cmdArgs.StdinData); err != nil {
		return err
	}

	switch cmd {
	case sp.ListReferrersCommand:
		return c.cmdListReferrers(cmdArgs, listReferrers)
	case sp.GetBlobContentCommand:
		return c.cmdGetBlob(cmdArgs, getBlobContent)
	case sp.GetRefManifestCommand:
		return c.cmdGetRefManifest(cmdArgs, getRefManifest)
	default:
		return plugin.NewError(types.ErrUnknownCommand, fmt.Sprintf("unknown %s: %v", sp.CommandEnvKey, cmd), "")
	}
}

func (c *pcontext) cmdListReferrers(cmdArgs *CmdArgs, pluginFunc ListReferrers) *plugin.Error {
	pluginArgs, err := plugin.ParseInputArgs(cmdArgs.Args)

	if err != nil {
		return plugin.NewError(types.ErrArgsParsingFailure, "error parsing args", err.Error())
	}

	var nextToken, artifactType string
	for _, arg := range pluginArgs {
		switch arg[0] {
		case "nextToken":
			nextToken = arg[1]
		case "artifactTypes":
			artifactType = arg[1]
		default:
			return plugin.NewError(types.ErrArgsParsingFailure, fmt.Sprintf("unknown args %s", arg[0]), "")
		}
	}

	result, err := pluginFunc(cmdArgs, cmdArgs.subjectRef, strings.Split(artifactType, ","), nextToken)

	if err != nil {
		return plugin.NewError(types.ErrPluginCmdFailure, fmt.Sprintf("plugin command %s failed", sp.ListReferrersCommand), err.Error())
	}

	err = types.WriteListReferrersResult(&result, c.Stdout)
	if err != nil {
		return plugin.NewError(types.ErrIOFailure, "failed to write plugin output", err.Error())
	}

	return nil
}

func (c *pcontext) cmdGetBlob(cmdArgs *CmdArgs, pluginFunc GetBlobContent) *plugin.Error {
	pluginArgs, err := plugin.ParseInputArgs(cmdArgs.Args)

	if err != nil {
		return plugin.NewError(types.ErrArgsParsingFailure, "error parsing args", err.Error())
	}

	var digestArg string
	for _, arg := range pluginArgs {
		switch arg[0] {
		case "digest":
			digestArg = arg[1]
		default:
			return plugin.NewError(types.ErrArgsParsingFailure, fmt.Sprintf("unknown args %s", arg[0]), "")
		}
	}

	digest, err := digest.Parse(digestArg)
	if err != nil {
		return plugin.NewError(types.ErrArgsParsingFailure, fmt.Sprintf("cannot parse digest arg %s", digestArg), err.Error())
	}

	result, err := pluginFunc(cmdArgs, cmdArgs.subjectRef, digest)

	if err != nil {
		return plugin.NewError(types.ErrPluginCmdFailure, fmt.Sprintf("plugin command %s failed", sp.ListReferrersCommand), err.Error())
	}

	_, err = c.Stdout.Write(result)
	if err != nil {
		return plugin.NewError(types.ErrIOFailure, "failed to write plugin output", err.Error())
	}

	return nil

}

func (c *pcontext) cmdGetRefManifest(cmdArgs *CmdArgs, pluginFunc GetReferenceManifest) *plugin.Error {
	pluginArgs, err := plugin.ParseInputArgs(cmdArgs.Args)

	if err != nil {
		return plugin.NewError(types.ErrArgsParsingFailure, "error parsing args", err.Error())
	}

	var digestArg string
	for _, arg := range pluginArgs {
		switch arg[0] {
		case "digest":
			digestArg = arg[1]
		default:
			return plugin.NewError(types.ErrArgsParsingFailure, fmt.Sprintf("unknown args %s", arg[0]), "")
		}
	}

	digest, err := digest.Parse(digestArg)
	if err != nil {
		return plugin.NewError(types.ErrArgsParsingFailure, fmt.Sprintf("cannot parse digest arg %s", digestArg), err.Error())
	}

	result, err := pluginFunc(cmdArgs, cmdArgs.subjectRef, digest)

	if err != nil {
		return plugin.NewError(types.ErrPluginCmdFailure, fmt.Sprintf("plugin command %s failed", sp.ListReferrersCommand), err.Error())
	}

	err = types.WriteReferenceManifestResult(&result, c.Stdout)
	if err != nil {
		return plugin.NewError(types.ErrIOFailure, "failed to write plugin output", err.Error())
	}

	return nil
}

func (c *pcontext) getCmdArgsFromEnv() (string, *CmdArgs, *plugin.Error) {
	argsMissing := make([]string, 0)

	// #1 Command
	var cmd = c.GetEnviron(sp.CommandEnvKey)
	if cmd == "" {
		argsMissing = append(argsMissing, sp.CommandEnvKey)
	}

	// #2 Version
	var version = c.GetEnviron(sp.VersionEnvKey)
	if version == "" {
		argsMissing = append(argsMissing, sp.VersionEnvKey)
	}

	// #3 Subject
	var subject = c.GetEnviron(sp.SubjectEnvKey)
	if subject == "" {
		argsMissing = append(argsMissing, sp.SubjectEnvKey)
	}

	// #4 Args
	var args = c.GetEnviron(sp.ArgsEnvKey)
	if args == "" {
		argsMissing = append(argsMissing, sp.ArgsEnvKey)
	}

	if len(argsMissing) > 0 {
		joined := strings.Join(argsMissing, ",")
		return "", nil, plugin.NewError(types.ErrMissingEnvironmentVariables, fmt.Sprintf("missing env variables [%s]", joined), "")
	}

	// TODO Limit the read length
	stdinData, err := ioutil.ReadAll(c.Stdin)
	if err != nil {
		return "", nil, plugin.NewError(types.ErrIOFailure, fmt.Sprintf("error reading from stdin: %v", err), "")
	}

	subRef, err := utils.ParseSubjectReference(subject)
	if err != nil {
		return "", nil, plugin.NewError(types.ErrArgsParsingFailure, fmt.Sprintf("cannot parse subject reference %s", subject), err.Error())
	}

	cmdArgs := &CmdArgs{
		Version:    version,
		Subject:    subject,
		Args:       args,
		StdinData:  stdinData,
		subjectRef: subRef,
	}

	return cmd, cmdArgs, nil
}

func validateVersion(version string, supportedVersions []string) *plugin.Error {
	for _, v := range supportedVersions {
		// TODO check for compatibility using semversion
		if v == version {
			return nil
		}
	}

	return plugin.NewError(types.ErrVersionNotSupported, fmt.Sprintf("plugin doesn't support version %s", version), "")
}

func validateConfig(jsonBytes []byte) *plugin.Error {
	var conf struct {
		Name string `json:"name"`
	}

	if err := json.Unmarshal(jsonBytes, &conf); err != nil {
		return plugin.NewError(types.ErrConfigParsingFailure, fmt.Sprintf("error unmarshall store config: %v", err), "")
	}
	if conf.Name == "" {
		return plugin.NewError(types.ErrInvalidStoreConfig, "missing store name", "")
	}
	return nil
}
