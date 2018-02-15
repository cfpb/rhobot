package healthcheck

import (
	"fmt"
	"strconv"
	"strings"
)

// TemplateHealthcheckPostgres pongo2 template for healthchecks INSERT
const TemplateHealthcheckPostgres = `
CREATE TABLE IF NOT EXISTS {{metadata.schema}}.{{metadata.table}}
(
  title text,
  query text,
  executed text,
  expected text,
  operation text,
  actual text,
  equal text,
  severity text,
  "timestamp" timestamp with time zone
);

INSERT INTO "{{metadata.schema}}"."{{metadata.table}}" ("title", "query", "executed", "expected", "operation", "actual", "equal", "severity", "timestamp") VALUES
{% for element in elements %}
('{{ element.Title }}', '{{ element.Query | safe | addquote }}', '{{ element.Passed}}', '{{ element.Expected  | safe | addquote  }}', '{{ element.Operation  | safe | addquote  }}', '{{ element.Actual  | safe | addquote  }}', '{{ element.Equal  | safe | addquote  }}', '{{ element.Severity }}', '{{ metadata.timestamp }}') ` +
	`{% if forloop.Last%};{%else%},{%endif%}` +
	`{% endfor %}`

// TemplateHealthcheckHTML pongo2 template for healthchecks
const TemplateHealthcheckHTML = `
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=us-ascii">
</head>
<style type="text/css">

body, p, h1, h3, ul, table {
		font-family: arial, sans-serif;
		font-size: 16px;
		color: #101820;
	}

	h1 {
		font-size: 34px;
		font-weight: normal;
	}

	h3 {
		font-size: 22px;
		font-weight: normal;
	}

	table {
		width: 100%;
		border-spacing: 0px;
	}

	td {
		text-align: left;
		padding: 8px;
	}

	td.entity {
		background-color: #addc91;
	}

	td.header_field {
		background-color: #e7e8e9;
		width: 13%;
	}

	td.data {
		border-bottom: 1px solid #b4b5b6;
	}

</style>



<h2>{{ metadata.status }}</h2>
<h2>{{ metadata.name }} - Running against database "{{ metadata.db_name }}"</h2>
<table>
	<tr>
		<td class = "header_field" >Title</td>
		<td class = "header_field" >Severity</td>
		<td class = "header_field" >Query</td>
		<td class = "header_field" >Test Ran?</td>
		<td class = "header_field" >Expected</td>
		<td class = "header_field" >Operation</td>
		<td class = "header_field" >Actual</td>
	</tr>
	{% for element in elements %}
	<tr>
		<td class = "data" >{{ element.Title }}</td>
		<td class = "data" >{{ element.Severity }}</td>
		<td class = "data" >{{ element.Query }}</td>

		{% if element.Equal == "TRUE"%}
			{% set bg_equals = "MediumSeaGreen" %}
		{% elif element.Equal == "FALSE" and element.Severity == "WARN"%}
			{% set bg_equals = "LightGoldenRodYellow" %}
		{% else %}
			{% set bg_equals = "LightCoral" %}
		{% endif %}

		<td class = "data"  bgcolor={{bg_equals}}>{{ element.Passed }}</td>
		{% if element.Passed == "SUCCESS"%}
		<td class = "data"  bgcolor={{bg_equals}}>{{ element.Expected }}</td>
		<td class = "data"  bgcolor={{bg_equals}}>{{ element.Operation }}</td>
		<td class = "data"  bgcolor={{bg_equals}}>{{ element.Actual }}</td>
		{% else %}
		<td class = "data"  bgcolor={{bg_equals}} colspan="3">{{ element.Error }}</td>
		{% endif %}

	{% endfor %}
	</tr>
</table>

{{ metadata.footer | safe }}<br> {{ metadata.timestamp }}
</html>

`

// FooterHealthcheck footer for healthchecks
const FooterHealthcheck = `
<p>Thank you,</p>
  <p>
  CFPB Data Team<br>
  Consumer Financial Protection Bureau
  </p>
  <p>Confidentiality Notice: If you received this email by mistake, please notify the sender of the mistake and delete the e-mail and any attachments. An inadvertent disclosure is not intended to waive any privileges.</p>
`

// SubjectHealthcheck creates a subject for healthcheck email
func SubjectHealthcheck(name string, dbName string, hostname string, level string, errors int, warnings int, fatal bool) string {

	hcName := name
	if name == "" {
		hcName = "healthchecks"
	}

	subjectStr := fmt.Sprintf("%s - %s - %s - %s level",
		hcName, dbName, hostname, strings.ToUpper(level))

	statusStr := StatusHealthchecks(errors, warnings, fatal)
	subjectStr = fmt.Sprintf("%s - %s", statusStr, subjectStr)

	return subjectStr
}

// StatusHealthchecks returns a simple summary for all healthchecks
func StatusHealthchecks(errors int, warnings int, fatal bool) string {

	if fatal {
		return fmt.Sprintf("FATAL")
	} else if errors > 0 {
		return fmt.Sprintf("ERROR(s) %s", strconv.Itoa(errors))
	} else if warnings > 0 {
		return fmt.Sprintf("WARNING(s) %s", strconv.Itoa(warnings))
	} else {
		return fmt.Sprintf("PASS")
	}
}
