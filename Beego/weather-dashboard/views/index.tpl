<!DOCTYPE html>
<html>
<head>
    <title>Weather Dashboard</title>

    <style>

        body{
            font-family: Arial;
            margin:40px;
        }

        table{
            border-collapse: collapse;
            width:400px;
        }

        th,td{
            border:1px solid #ccc;
            padding:10px;
        }

    </style>
</head>
<body>

<h1>Weather Dashboard</h1>

<table>

    <tr>
        <th>City</th>
        <th>Temperature</th>
    </tr>

    {{range .Weather}}

    <tr>
        <td>{{.City}}</td>
        <td>{{.Temperature}} °C</td>
    </tr>

    {{end}}

</table>

</body>
</html>