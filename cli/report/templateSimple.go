package report

// TemplateSimple constant for simple html pongo2 template
const TemplateSimple = `
<h2>Test type: {{ metadata.test }}</h2>
{% for element in elements %}
 <tr>
   <td>Some: {{ element.Some }} </td>
   <td>Thing: {{ element.Thing }}</td>
</tr>
{% endfor %}
`
