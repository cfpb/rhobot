package healthcheck

// TemplateHealthcheck pongo2 template for healthchecks
const TemplateHealthcheck = `
  <h2>{{ metadata.name }} - Running against database "{{ metadata.db_name }}"</h2>
    <table border=1 frame=void rules=rows>
        <tr>
            <th>Title</th>
            <th>Query</th>
            <th>Test Ran?</th>
            <th>Expected</th>
            <th>Actual</th>
        </tr>
        {% for element in elements %}
            <tr>
                <td>{{ element.Title }}</td>
                <td>{{ element.Query }}</td>
                {% if element.Passed == "SUCCESS"%}
                    <td bgcolor="green">{{ element.Passed }}</td>
                    <td bgcolor="green">{{ element.Expected }}</td>
                    <td bgcolor="green">{{ element.Actual }}</td>
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
  CFPB CR Data-Sharing Team<br>
  Consumer Financial Protection Bureau
  </p>
  <p>Confidentiality Notice: If you received this email by mistake, please notify the sender of the mistake and delete the e-mail and any attachments. An inadvertent disclosure is not intended to waive any privileges.</p>
`
