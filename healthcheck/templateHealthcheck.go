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
  actual text,
  severity text,
  "timestamp" timestamp with time zone
);

INSERT INTO "{{metadata.schema}}"."{{metadata.table}}" ("title", "query", "executed", "expected", "actual", "severity", "timestamp") VALUES
{% for element in elements %}
('{{ element.Title }}', '{{ element.Query | safe  }}', '{{ element.Passed}}', '{{ element.Expected }}', '{{ element.Actual }}', '{{ element.Severity }}', '{{ metadata.timestamp }}') ` +
	`{% if forloop.Last%};{%else%},{%endif%}` +
	`{% endfor %}`

// TemplateHealthcheckHTML pongo2 template for healthchecks
const TemplateHealthcheckHTML = `
	<h2>{{ metadata.status }}</h2>
  <h2>{{ metadata.name }} - Running against database "{{ metadata.db_name }}"</h2>
    <table border=1 frame=void rules=rows>
        <tr>
            <th>Title</th>
						<th>Severity</th>
            <th>Query</th>
            <th>Test Ran?</th>
            <th>Expected</th>
            <th>Actual</th>
        </tr>
        {% for element in elements %}
            <tr>
                <td>{{ element.Title }}</td>
								<td>{{ element.Severity }}</td>
                <td>{{ element.Query }}</td>
                {% if element.Passed == "SUCCESS"%}
                    <td bgcolor="green">{{ element.Passed }}</td>
										{% if element.Equal == "TRUE"%}
                    	<td bgcolor="green">{{ element.Expected }}</td>
                    	<td bgcolor="green">{{ element.Actual }}</td>
										{% elif  element.Equal == "FALSE" %}
                    	<td bgcolor="red">{{ element.Expected }}</td>
                    	<td bgcolor="red">{{ element.Actual }}</td>
										{% endif %}
                {% elif  element.Passed == "FAIL" %}
                    <td bgcolor="red">{{ element.Passed }}</td>
                    <td>{{ element.Error }}</td>
                {% endif %}
            </tr>
        {% endfor %}
    </table>
    {{ metadata.footer | safe }}<br>
    {{ metadata.timestamp }}

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
func SubjectHealthcheck(name string, dbName string, hostname string, level string, errors int, fatal bool) string {

	hcName := name
	if name == "" {
		hcName = "healthchecks"
	}

	subjectStr := fmt.Sprintf("%s - %s - %s - %s level",
		hcName, dbName, hostname, strings.ToUpper(level))

	statusStr := StatusHealthchecks(errors, fatal)
	subjectStr = fmt.Sprintf("%s - %s", statusStr, subjectStr)

	return subjectStr
}

// StatusHealthchecks returns a simple summray for all healthchecks
func StatusHealthchecks(errors int, fatal bool) string {

	if fatal {
		return fmt.Sprintf("FATAL")
	} else if errors > 0 {
		return fmt.Sprintf("ERROR(s) %s", strconv.Itoa(errors))
	} else {
		return fmt.Sprintf("PASS")
	}
}
