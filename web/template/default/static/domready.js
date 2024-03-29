$(document).ready(function() {
    let loaderHTML = "<div class=\"spinner-grow text-primary\" role=\"status\"><span class=\"sr-only\">Loading...</span></div>" +
        "<div class=\"spinner-grow text-secondary\" role=\"status\"><span class=\"sr-only\">Loading...</span></div>" +
        "<div class=\"spinner-grow text-success\" role=\"status\"><span class=\"sr-only\">Loading...</span></div>" +
        "<div class=\"spinner-grow text-danger\" role=\"status\"><span class=\"sr-only\">Loading...</span></div>" +
        "<div class=\"spinner-grow text-warning\" role=\"status\"><span class=\"sr-only\">Loading...</span></div>" +
        "<div class=\"spinner-grow text-info\" role=\"status\"><span class=\"sr-only\">Loading...</span></div>" +
        "<div class=\"spinner-grow text-light\" role=\"status\"><span class=\"sr-only\">Loading...</span></div>" +
        "<div class=\"spinner-grow text-dark\" role=\"status\"><span class=\"sr-only\">Loading...</span></div>";

    $('#aboutModal').on('show.bs.modal', function (event) {
        let button = $(event.relatedTarget); // Button that triggered the modal
        let aboutDataURL = button.data('url'); // Extract info from data-* attributes
        let appName = button.data('appname');
        // If necessary, you could initiate an AJAX request here (and then do the updating in a callback).
        // Update the modal's content. We'll use jQuery here, but you could use a data binding library or other methods instead.
        var modal = $(this);
        $.ajax({
            url: aboutDataURL,
            beforeSend: function(jqXHR, settings) {
                modal.find('.modal-title').text('About '+appName);
                modal.find('.modal-body').html(loaderHTML);
            },
            success: function (data, textStatus, jqXHR) {
                console.log("aboutModal AJAX success. data = "+JSON.stringify(data));
                htmlData = "<table class=\"table text-white\">" +
                    "<tr><th>Version</th><td>" + data.version + "</td></tr>" +
                    "<tr><th>Build Data</th><td>" + data.build_date + "</td><tr>" +
                    "<tr><th>Git Commit</th><td>" + data.git_version + "</td></tr>" +
                    "<tr><th>Git Branch</th><td>" + data.git_branch + "</td></tr>" +
                    "<tr><th>Git State</th><td>" + data.git_state + "</td></tr>" +
                    "<tr><th>Git Summary</th><td>" + data.git_summary + "</td></tr>" +
                    "</table>";
                modal.find('.modal-body').html(htmlData);
            },
            error: function(jqXHR, textStatus, errorThrown) {
                console.log("jqXHR = ", JSON.stringify(jqXHR));
                console.log("textStatus = ", textStatus);
                console.log("errorThrown = ", errorThrown);
                htmlData = "<span class=\"card bg-light p-3 text-danger\"><p class=\"mb-0 text-sc\">Unable to fetch About Information. Got: "+JSON.stringify(jqXHR)+"</p></span>";
                modal.find('.modal-body').html(htmlData);
            },
            timeout: 10000
        });
    });

    if( $("#cron-list-data").length ) {
        console.log("cron-list-data.length is true. Refresh event will be scheduled");
        dataURL = $("#cron-list-data").data("url");
        refreshInterval = $("#cron-list-data").data("refresh");
        // Refresh the cron table every 5 seconds
        setInterval(function () {
            $.get(dataURL, function (crons) {
                console.log("crons data = ", JSON.stringify(crons))
                // Clear the cron table
                $("#cron-table tbody").empty();
                $("#cron-list-data").html(loaderHTML);
                // Populate the cron table
                let cronListTableData = "<table class=\"table table-bordered table-white\" id=\"cron-list-table\" style=\"color:white\"><thead><tr>" +
                    "<th>ID</th><th>Schedule</th><th>Next Execution</th><th>Previous Execution</th></tr></thead><tbody>";
                crons.forEach(function (cron) {
                    console.log("cron data = ", JSON.stringify(cron));
                    cronListTableData += "<tr><td>" + cron.jobID + "</td><td>" + cron.schedule + "</td><td>" + cron.next + "</td><td>" + cron.previous + "</td></tr>";
                });

                cronListTableData += "</tbody></table>";
                $("#cron-list-data").html(cronListTableData);
            });
        }, refreshInterval);
    } else {
        console.log("cron-list-data.length is false. Refresh event will not be scheduled");
    }
});