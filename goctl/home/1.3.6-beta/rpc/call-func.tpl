
{{if .hasComment}}{{.comment}}{{end}}
func (m *default{{.serviceName}}) {{.method}}(ctx context.Context{{if .hasReq}}, in *{{.pbRequest}}{{end}}, opts ...grpc.CallOption) ({{if .notStream}}*{{.pbResponse}}, {{else}}{{.streamBody}},{{end}} error) {
	client := {{if .isCallPkgSameToGrpcPkg}}{{else}}{{.package}}.{{end}}New{{.rpcServiceName}}Client(m.cli.Conn())
	return client.{{.method}}(ctx{{if .hasReq}}, in{{end}}, opts...)
}
