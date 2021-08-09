package runtimes

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"path/filepath"

	"github.com/concourse/go-archive/tarfs"
	"github.com/vito/bass"
	"gopkg.in/yaml.v2"
)

func NewEnv(cwd string, runtime Runtime) *bass.Env {
	env := bass.NewStandardEnv()

	env.Set("run",
		bass.Func("run", func(ctx context.Context, workload bass.Workload) (*bass.Source, error) {
			err := runtime.Run(ctx, workload)
			if err != nil {
				return nil, err
			}

			buf := new(bytes.Buffer)
			err = runtime.Response(ctx, buf, workload)
			if err != nil {
				return nil, err
			}

			name, err := workload.SHA1()
			if err != nil {
				return nil, err
			}

			return bass.NewSource(bass.NewJSONSource(name, buf)), nil
		}),
		`run a workload`)

	env.Set("path",
		bass.Func("path", func(workload bass.Workload, path bass.FileOrDirPath) bass.WorkloadPath {
			return bass.WorkloadPath{
				Workload: workload,
				Path:     path,
			}
		}),
		`returns a path within a workload`)

	env.Set("export",
		bass.Func("export", func(ctx context.Context, path bass.WorkloadPath, dest bass.DirPath) error {
			r, w := io.Pipe()
			go func() {
				w.CloseWithError(runtime.Export(ctx, w, path.Workload, path.Path.FilesystemPath()))
			}()

			return tarfs.Extract(r, filepath.Join(cwd, dest.FromSlash()))
		}),
		`export a workload path to a local path`)

	env.Set("yaml-decode",
		bass.Func("yaml-decode", func(ctx context.Context, path bass.WorkloadPath) (bass.Value, error) {
			buf := new(bytes.Buffer)
			err := runtime.Export(ctx, buf, path.Workload, path.Path.FilesystemPath())
			if err != nil {
				return nil, err
			}

			tr := tar.NewReader(buf)
			_, err = tr.Next()
			if err != nil {
				return nil, err
			}

			var v interface{}
			err = yaml.NewDecoder(tr).Decode(&v)
			if err != nil {
				return nil, err
			}

			return bass.ValueOf(v)
		}))

	return bass.NewEnv(env)
}
