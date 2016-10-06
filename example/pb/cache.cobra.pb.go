// Code generated by protoc-gen-cobra.
// source: cache.proto
// DO NOT EDIT!

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	io "io"
	json "encoding/json"
	log "log"
	cobra "github.com/spf13/cobra"
	filepath "path/filepath"
	ioutil "io/ioutil"
	time "time"
	tls "crypto/tls"
	os "os"
	pflag "github.com/spf13/pflag"
	template "text/template"
	x509 "crypto/x509"
	context "golang.org/x/net/context"
	credentials "google.golang.org/grpc/credentials"
	grpc "google.golang.org/grpc"
	iocodec "github.com/fiorix/protoc-gen-cobra/iocodec"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ cobra.Command
var _ filepath.WalkFunc
var _ = ioutil.Discard
var _ time.Time
var _ tls.Config
var _ os.File
var _ pflag.FlagSet
var _ template.Template
var _ x509.Certificate
var _ context.Context
var _ credentials.AuthInfo
var _ grpc.ClientConn
var _ iocodec.Encoder
var _ io.Reader
var _ json.Encoder
var _ log.Logger

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

var _DefaultCacheClientCommandConfig = _NewCacheClientCommandConfig()

func init() {
	_DefaultCacheClientCommandConfig.AddFlags(CacheClientCommand.PersistentFlags())
}

type _CacheClientCommandConfig struct {
	ServerAddr         string
	RequestFile        string
	PrintSampleRequest bool
	ResponseFormat     string
	Timeout            time.Duration
	TLS                bool
	InsecureSkipVerify bool
	CACertFile         string
	CertFile           string
	KeyFile            string
}

func _NewCacheClientCommandConfig() *_CacheClientCommandConfig {
	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = "localhost:8080"
	}
	timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		timeout = 10 * time.Second
	}
	outfmt := os.Getenv("RESPONSE_FORMAT")
	if outfmt == "" {
		outfmt = "json"
	}
	return &_CacheClientCommandConfig{
		ServerAddr:         addr,
		RequestFile:        os.Getenv("REQUEST_FILE"),
		PrintSampleRequest: os.Getenv("PRINT_SAMPLE_REQUEST") != "",
		ResponseFormat:     outfmt,
		Timeout:            timeout,
		TLS:                os.Getenv("TLS") != "",
		InsecureSkipVerify: os.Getenv("TLS_INSECURE_SKIP_VERIFY") != "",
		CACertFile:         os.Getenv("TLS_CA_CERT_FILE"),
		CertFile:           os.Getenv("TLS_CERT_FILE"),
		KeyFile:            os.Getenv("TLS_KEY_FILE"),
	}
}

func (o *_CacheClientCommandConfig) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.ServerAddr, "server-addr", "s", o.ServerAddr, "server address in form of host:port")
	fs.StringVarP(&o.RequestFile, "request-file", "f", o.RequestFile, "client request file (must be json, yaml, or xml); use \"-\" for stdin + json")
	fs.BoolVarP(&o.PrintSampleRequest, "print-sample-request", "p", o.PrintSampleRequest, "print sample request file and exit")
	fs.StringVarP(&o.ResponseFormat, "response-format", "o", o.ResponseFormat, "response format (json, prettyjson, yaml, or xml)")
	fs.DurationVar(&o.Timeout, "timeout", o.Timeout, "client connection timeout")
	fs.BoolVar(&o.TLS, "tls", o.TLS, "enable tls")
	fs.BoolVar(&o.InsecureSkipVerify, "tls-insecure-skip-verify", o.InsecureSkipVerify, "INSECURE: skip tls checks")
	fs.StringVar(&o.CACertFile, "tls-ca-cert-file", o.CACertFile, "ca certificate file")
	fs.StringVar(&o.CertFile, "tls-cert-file", o.CertFile, "client certificate file")
	fs.StringVar(&o.KeyFile, "tls-key-file", o.KeyFile, "client key file")
}

var CacheClientCommand = &cobra.Command{
	Use: "cache",
}

func _DialCache() (*grpc.ClientConn, CacheClient, error) {
	cfg := _DefaultCacheClientCommandConfig
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(cfg.Timeout),
	}
	if cfg.TLS {
		tlsConfig := tls.Config{}
		if cfg.InsecureSkipVerify {
			tlsConfig.InsecureSkipVerify = true
		}
		if cfg.CACertFile != "" {
			cacert, err := ioutil.ReadFile(cfg.CACertFile)
			if err != nil {
				return nil, nil, fmt.Errorf("ca cert: %v", err)
			}
			certpool := x509.NewCertPool()
			certpool.AppendCertsFromPEM(cacert)
			tlsConfig.RootCAs = certpool
		}
		if cfg.CertFile != "" {
			if cfg.KeyFile == "" {
				return nil, nil, fmt.Errorf("missing key file")
			}
			pair, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
			if err != nil {
				return nil, nil, fmt.Errorf("cert/key: %v", err)
			}
			tlsConfig.Certificates = []tls.Certificate{pair}
		}
		cred := credentials.NewTLS(&tlsConfig)
		opts = append(opts, grpc.WithTransportCredentials(cred))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(cfg.ServerAddr, opts...)
	if err != nil {
		return nil, nil, err
	}
	return conn, NewCacheClient(conn), nil
}

type _CacheRoundTripFunc func(cli CacheClient, in iocodec.Decoder, out iocodec.Encoder) error

func _CacheRoundTrip(sample interface{}, fn _CacheRoundTripFunc) error {
	cfg := _DefaultCacheClientCommandConfig
	var em iocodec.EncoderMaker
	var ok bool
	if cfg.ResponseFormat == "" {
		em = iocodec.DefaultEncoders["json"]
	} else {
		em, ok = iocodec.DefaultEncoders[cfg.ResponseFormat]
		if !ok {
			return fmt.Errorf("invalid response format: %q", cfg.ResponseFormat)
		}
	}
	if cfg.PrintSampleRequest {
		return em.NewEncoder(os.Stdout).Encode(sample)
	}
	var d iocodec.Decoder
	if cfg.RequestFile == "" || cfg.RequestFile == "-" {
		d = iocodec.DefaultDecoders["json"].NewDecoder(os.Stdin)
	} else {
		f, err := os.Open(cfg.RequestFile)
		if err != nil {
			return fmt.Errorf("request file: %v", err)
		}
		defer f.Close()
		ext := filepath.Ext(cfg.RequestFile)
		if len(ext) > 0 && ext[0] == '.' {
			ext = ext[1:]
		}
		dm, ok := iocodec.DefaultDecoders[ext]
		if !ok {
			return fmt.Errorf("invalid request file format: %q", ext)
		}
		d = dm.NewDecoder(f)
	}
	conn, client, err := _DialCache()
	if err != nil {
		return err
	}
	defer conn.Close()
	return fn(client, d, em.NewEncoder(os.Stdout))
}

var _CacheSetClientCommand = &cobra.Command{
	Use: "set",
	Run: func(cmd *cobra.Command, args []string) {
		var v SetRequest
		err := _CacheRoundTrip(v, func(cli CacheClient, in iocodec.Decoder, out iocodec.Encoder) error {

			err := in.Decode(&v)
			if err != nil {
				return err
			}

			resp, err := cli.Set(context.Background(), &v)

			if err != nil {
				return err
			}

			return out.Encode(resp)

		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	CacheClientCommand.AddCommand(_CacheSetClientCommand)
}

var _CacheGetClientCommand = &cobra.Command{
	Use: "get",
	Run: func(cmd *cobra.Command, args []string) {
		var v GetRequest
		err := _CacheRoundTrip(v, func(cli CacheClient, in iocodec.Decoder, out iocodec.Encoder) error {

			err := in.Decode(&v)
			if err != nil {
				return err
			}

			resp, err := cli.Get(context.Background(), &v)

			if err != nil {
				return err
			}

			return out.Encode(resp)

		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	CacheClientCommand.AddCommand(_CacheGetClientCommand)
}

var _CacheMultiSetClientCommand = &cobra.Command{
	Use: "multiset",
	Run: func(cmd *cobra.Command, args []string) {
		var v SetRequest
		err := _CacheRoundTrip(v, func(cli CacheClient, in iocodec.Decoder, out iocodec.Encoder) error {

			stream, err := cli.MultiSet(context.Background())
			if err != nil {
				return err
			}
			for {
				err = in.Decode(&v)
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}
				err = stream.Send(&v)
				if err != nil {
					return err
				}
			}

			resp, err := stream.CloseAndRecv()
			if err != nil {
				return err
			}

			return out.Encode(resp)

		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	CacheClientCommand.AddCommand(_CacheMultiSetClientCommand)
}

var _CacheMultiGetClientCommand = &cobra.Command{
	Use: "multiget",
	Run: func(cmd *cobra.Command, args []string) {
		var v GetRequest
		err := _CacheRoundTrip(v, func(cli CacheClient, in iocodec.Decoder, out iocodec.Encoder) error {

			stream, err := cli.MultiGet(context.Background())
			if err != nil {
				return err
			}
			for {
				err = in.Decode(&v)
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}
				err = stream.Send(&v)
				if err != nil {
					return err
				}
			}

			for {
				v, err := stream.Recv()
				if err != nil {
					return err
				}
				err = out.Encode(v)
				if err != nil {
					return err
				}
			}
			return nil

		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	CacheClientCommand.AddCommand(_CacheMultiGetClientCommand)
}
