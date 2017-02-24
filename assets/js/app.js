function lexical(a, b) {
    return a < b ? -1
         : a > b ?  1
         : 0;
}
function oldest(list, field) {
    var x = list[0][field];
    for (var i = 1; i < list.length; i++) {
        if (list[i][field] < x) {
            x = list[i][field];
        }
    }
    return x;
}

function css(ts) {
    var s = (new Date()).getTime() / 1000 - ts;
    if (s <   5 * 86400) { return 'age ok' }
    if (s <  90 * 86400) { return 'age warn' }
    return 'age crit';
}

function age(ts) {
    var s = (new Date()).getTime() / 1000 - ts;
    s = parseInt(s / 86400);
    if (s <   1) { return 'new'  }
    if (s <  31) { return s.toString() + 'd' }
    return parseInt(s / 7).toString() + 'w';
}

function bucket(ts) {
    var s = (new Date()).getTime() / 1000 - ts;
    s = parseInt(s / 86400);
    if (s <=   1) { return 9; } /* 1d */
    if (s <=   2) { return 8; } /* 2d */
    if (s <=   3) { return 7; } /* 3d */
    if (s <=   7) { return 6; } /* 1w */
    if (s <=  30) { return 5; } /* 1m */
    if (s <=  60) { return 4; } /* 2m */
    if (s <=  90) { return 3; } /* 3m */
    if (s <= 180) { return 2; } /* 6m */
    if (s <= 270) { return 1; } /* 9m */
                    return 0;   /*  + */
}

function histogram(db, type) {
    var buckets = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0];
    for (var k in db) {
        for (var i = 0; i < db[k][type].length; i++) {
            buckets[bucket(db[k][type][i].updated)]++;
        }
    }

    var names = ['older', '9mo', '6mo', '3mo', '2mo', '1mo', '1w', '3d', '2d', 'new'];
    for (var i = 0; i < buckets.length; i++) {
        buckets[i] = {
            bucket: names[i],
            count:  buckets[i]
        };
    }
    return buckets;
}

function histograph(data, svg) {
  svg.html("");
  var margin = {top: 0, right: 0, bottom: 20, left: 0},
      width  = svg.attr("width")  - margin.left - margin.right,
      height = svg.attr("height") - margin.top  - margin.bottom;

  var g = svg.append("g")
             .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

  var x = d3.scaleBand().rangeRound([0, width]).padding(0.1),
      y = d3.scaleLinear().rangeRound([height, 0]);

  x.domain(data.map(function (d) { return d.bucket; }));
  y.domain([0, d3.max(data, function (d) { return d.count })]);

  g.append("g")
   .attr("class", "axis axis--x")
   .attr("transform", "translate(0," + height + ")")
   .call(d3.axisBottom(x));

  g.selectAll('bar')
   .data(data)
   .enter().append("rect")
     .attr("class", "bar")
     .attr("x", function(d) { return x(d.bucket); })
     .attr("y", function(d) { return y(d.count); })
     .attr("width", x.bandwidth())
     .attr("height", function(d) { return height - y(d.count); });
}

function assignedto(thing) {
    if (thing.assignees.length == 0) {
        return '<em>unassigned</em>';
    }
    return '<strong>'+thing.assignees.join('</strong>, <strong>')+'</strong>';
}

$(function () {
    $('#dashboard, #configure').hide();

    var showDashboard = function() {
        $('#configure').hide();
        $.ajax({
            type: 'GET',
            url:  '/v1/health',
            success: function (data) {

                /* re-order based on oldest issue/pr */
                var pulls  = [];
                var issues = [];
                var $l;

                for (var k in data) {
                    if (data[k].pulls.length > 0) {
                        pulls.push([
                            oldest(data[k].pulls, 'updated'),
                            data[k]
                        ]);
                    }
                    if (data[k].issues.length > 0) {
                        issues.push([
                            oldest(data[k].issues, 'updated'),
                            data[k]
                        ]);
                    }
                }
                issues = issues.sort(function (a, b) { return a[0] - b[0] });
                pulls  =  pulls.sort(function (a, b) { return a[0] - b[0] });

                $l = $('#issues .list').empty();
                for (var i = 0; i < issues.length; i++) {
                    var list = '';
                    var repo = issues[i][1];
                    repo.issues = repo.issues.sort(function (a, b) { return a.updated - b.updated });
                    for (var j = 0; j < repo.issues.length; j++) {
                        var issue = repo.issues[j];
                        list += '<li><a href="'+issue.url+'" target="_blank">'+
                                       '#'+issue.number+' '+issue.title+'</a>'+
                                    '<span class="'+css(issue.updated)+'">'+age(issue.updated)+'</span>'+
                                    '<p>'+assignedto(issue)+'</p></li>';
                    }
                    $l.append($('<h3>'+repo.org+' / '+repo.name+'</h3>'+
                                '<ul>'+list+'</ul>'));
                }
                histograph(histogram(data, 'issues'), d3.select('#issues svg'));

                $l = $('#pulls .list').empty();
                for (var i = 0; i < pulls.length; i++) {
                    var list = '';
                    var repo = pulls[i][1];
                    repo.pulls = repo.pulls.sort(function (a, b) { return a.updated - b.updated });
                    for (var j = 0; j < repo.pulls.length; j++) {
                        var pull = repo.pulls[j];
                        list += '<li><a href="'+pull.url+'" target="_blank">'+
                                       '#'+pull.number+' '+pull.title+'</a>'+
                                    '<span class="'+css(pull.updated)+'">'+age(pull.updated)+'</span>'+
                                    '<p>'+assignedto(pull)+'</p></li>';
                    }
                    $l.append($('<h3>'+repo.org+' / '+repo.name+'</h3>'+
                                '<ul>'+list+'</ul>'));
                }
                histograph(histogram(data, 'pulls'), d3.select('#pulls svg'));

                $('#dashboard').show();
            },
        });
    };

    var showConfigure = function() {
        $('#dashboard').hide();
        $.ajax({
            type: 'GET',
            url:  '/v1/repos',
            success: function (data) {
                data = data.sort(function (a, b) {
                    if (a.org < b.org) return -1;
                    if (a.org > b.org) return  1;
                    if (a.name < b.name) return -1;
                    if (a.name < b.name) return  1;
                    return 0;
                });
                var $ul = $('#configure ul').empty();
                for (var i = 0; i < data.length; i++) {
                    var checked = (data[i].included ? 'checked="checked" ' : '');
                    $ul.append($('<li><input type="checkbox" '+checked+'name="'+data[i].id+'">'+
                                      data[i].org + ' / ' + data[i].name + '</li>'));
                }
                $('#configure').show();
            }
        });
    };



    $(document.body).on('click', 'a[href="#config"]', function (event) {
        event.preventDefault();
        showConfigure();

    }).on('click', 'a[href="#home"]', function (event) {
        event.preventDefault();
        showDashboard();

    }).on('change', '#configure input[type=checkbox]', function (event) {
        $.ajax({
            type: 'POST',
            url:  '/v1/repos',
            processData: false,
            data: JSON.stringify($('#configure').serializeArray()),
            success: function () {
                console.log('ok');
            }
        });

    }).on('click', 'button[rel=scrape]', function (event) {
        event.preventDefault();
        $.ajax({
            type: 'POST',
            url:  '/v1/scrape',
            processData: false,
            data: JSON.stringify($('#configure').serializeArray()),
            success: function () {
                console.log('scraped');
            }
        });

    });
    showDashboard();
});
