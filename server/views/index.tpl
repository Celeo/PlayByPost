{{define "content"}}

<!-- TODO pagination -->

{{range .posts}}
  <div class="card rounded shadow">
    <div class="card-body">
      <h5 class="card-title clearfix">
        <span class="float-left">{{.Name}} ({{.Tag}})</span>
        <span class="float-right">{{.Date}}</span>
      </h5>
      {{.Content}}
      <div id="rolls">
        {{range .Rolls}}
          <div class="roll">
            {{.String}} => {{.Value}}{{if .IsD20Crit}} (crit){{end}}
          </div>
        {{end}}
      </div>
    </div>
  </div>
  <br>
{{end}}

<!-- TODO new post widget -->

{{end}}
