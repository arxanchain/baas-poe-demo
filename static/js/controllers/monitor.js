/* global ace_config */
/* global $ */
var bag = {
	peers: [],
	chain: {},
	stats:{
		height: 0,
		oldest_block: 999,
		session_blocks: 0,
		trans: 0,
		deploys: 0,
		invokes: 0,
		lastHr: 0
	}};
var get_last = 10;
var known_blocks = {};

// =================================================================================
// On Load
// =================================================================================
$(document).on('ready', function() {

	// 1. get peer/user list
	rest_get_audit(cb_got_audit);

	// 2. register user
	function cb_got_audit(e, resp){
		rest_register(resp.user.username, resp.user.secret, cb_registered);
	}

	// 3. get chain stats
	function cb_registered(e, resp){
		rest_chainstats(cb_got_chainstats);
	}

	// 4. get block stats
	var next = 0;
	function cb_got_chainstats(e, resp){
		if(resp.height > bag.stats.height){
			console.log('new block height!');
			bag.chain = resp;
			if(next == 0) next = resp.height - get_last + 1;			//only get the last X if its the first time around
			if(next < 0) next = 0;

			rest_blockstats(next, get_blocks);
		}
	}

	// 5-n. continue gettting block stats
	function get_blocks(){
		next++;
		if(next <= bag.chain.height){
			rest_blockstats(next, get_blocks);
		}
	}

	// Periodically check on height
	setInterval(function(){
		rest_chainstats(cb_got_chainstats);
		$(".dateText").each(function(){
			var i = $(this).attr('height');
			if(known_blocks[i] && known_blocks[i].transactions > 0){
				$(this).html(formatTime(known_blocks[i].transactions[0].timestamp.seconds) + " <br/>ago");
			}
		});
	}, 10000);


	// =================================================================================
	// jQuery UI Events
	// =================================================================================
	$(document).on("click", "#loadMore", function(){
		console.log('starting at', bag.stats.oldest_block);
		last = bag.stats.oldest_block;
		count = 0;
		get_prev_blocks();
		return false;
	});

	$(document).on("click", ".block", function(){
		var height = $(this).attr('height');
		build_row(height);
		$(".selectedBlock").removeClass("selectedBlock");

		$(".blockWrap").animate({left: "0"}, 0, function(){

		});

		$(this).parent().animate({left: "-=70"}, 300, function(){		//wiggle wiggle
			$(this).animate({left: "+=10"}, 300, function(){
				$(this).animate({left: "-=6"}, 300, function(){
				});
			});
		});
	});
});

var last = 0;
var count = 0;
var goingDown = false;
function get_prev_blocks(){
	last--;
	count++;
	if(count < get_last){
		if(last > 0){
			goingDown = true;
			rest_blockstats(last, get_prev_blocks);
		}
		else {
			goingDown = false;
			$("#loadMore").parent().parent().remove();
		}
	}
	else {
		goingDown = false;
		$("#loadMore").parent().parent().remove();
		$("#chainWrap").append('<div class="blockWrap">' +
										'<div class="dateWrap"></div>' +
										'<div class="block">' +
											'<div id="loadMore" class="height">load more</div>' +
										'</div>' +
									'</div>');
	}
}


// =================================================================================
// Block Fun
// =================================================================================
//receive new block
function new_block(block){
	if(!known_blocks[block.height]){
		add2stats(block);
		build_block(block);
	}

	known_blocks[block.height] = block;
}

//build UI block
function build_block(block){
	var html = '';
	var time = '';
	var deploys = 0;
	var invokes = 0;
	if(block.transactions && block.transactions[0]) {
		time = formatTime(block.transactions[0].timestamp.seconds);
		for(var i in block.transactions){
			if(block.transactions[i].type === 1) deploys++;
			if(block.transactions[i].type === 3) invokes++;
		}
	}

	html += '<div class="blockWrap" style="display:none;">';
	html +=		'<div class="dateWrap">';
	html +=			'<div class="date">';
	html +=				'<div class="dateText" height="' + block.height + '">' + time + '<br/>ago</div>';
	html +=			'</div>';
	html +=			'<div class="bar">&nbsp;</div>';
	html +=		'</div>';
	html +=		'<div class="block" height=' + block.height + '>';
	html +=			'<div class="height">' + block.height + '</div>';
	html +=			'<div class="deploy">' + deploys +' Deployment(s)</div>';
	html +=			'<div class="invoke">' + invokes + ' Invocation(s)</div>';
	html +=		'</div>';
	html +=	'</div>';

	if(goingDown){
		$("#chainWrap").append(html);
		$(".blockWrap").fadeIn();
	}
	else{
		$("#chainWrap").animate({top: "+=70"}, 400, function(){
			$("#chainWrap").prepend(html).css("top", "0");
			$(".blockWrap").fadeIn();
		});
	}
}

//build table row
function build_row(height){
	var html = '';

	for(var i in known_blocks[height].transactions){
		var ccid = atob(known_blocks[height].transactions[i].chaincodeID);
		var payload = atob(known_blocks[height].transactions[i].payload);
		var pos = payload.indexOf(ccid);
		var uuid = known_blocks[height].transactions[i].uuid;
		payload = payload.substring(pos + ccid.length + 2);
		if(known_blocks[height].transactions[i].type == 1) {				//if its a deploy, switch uuid and ccid for some unkown reason
			uuid = 'n/a';
			ccid = known_blocks[height].transactions[i].uuid;
		}
		ccid = ccid.substring(0, 12) + '...';

		html += '<tr>';
		html +=		'<td>' + formatDate(known_blocks[height].transactions[i].timestamp.seconds * 1000, '%M/%d %I:%m%p') + '</td>';
		html += 	'<td>' + type2word(known_blocks[height].transactions[i].type) + '</td>';
		html += 	'<td>' + uuid + '</td>';
		html += 	'<td>' + ccid + '</td>';
		html += 	'<td>' + payload + '</td>';
		html += '</tr>';
	}
	$("#detailsBody").html(html);
}

function type2word(type){
	if(type === 1) return 'DEPLOY';
	if(type === 3) return 'INVOKE';
	return type;
}

//record statistics
function add2stats(block){
	if(block.height && block.transactions){
		bag.stats.session_blocks++;
		if(block.height > bag.stats.height) bag.stats.height = block.height;
		if(block.height < bag.stats.oldest_block) bag.stats.oldest_block = block.height;

		var elasped = Date.now()/1000  - block.transactions[0].timestamp.seconds;

		if(elasped < 60*60*1) bag.stats.lastHr++;
		for(var i in block.transactions){
			if(block.transactions[i].type === 1) bag.stats.deploys++;
			if(block.transactions[i].type === 3) bag.stats.invokes++;
		}

		bag.stats.transPerBlock + block.transactions.length;

		$("#blockHeight").html(bag.stats.height);
		$("#blockDeploys").html(bag.stats.deploys);
		$("#blockInvokes").html(bag.stats.invokes);
		$("#blockRate").html(bag.stats.lastHr);
		$("#blockTrans").html(((bag.stats.deploys + bag.stats.invokes) / bag.stats.session_blocks).toFixed(1));
		$(".sessionBlocks").html(bag.stats.session_blocks);
	}
}


// =================================================================================
// REST fun
// =================================================================================
//rest call to service to get peers/users
function rest_get_audit(cb){
	//console.log("audit");
	//var host = 'broker-dev.obchain';
	//var port ='3000';

	$.ajax({
		method: 'GET',
		url: window.location.origin + '/api/audit/' + ace_config.network_id,
		contentType: 'application/json',
		success: function(json){
			console.log('Success - audit', json);
			bag.peers.push(json.peer);
			if(cb) cb(null, json);
		},
		error: function(e){
			console.log('Error - audit', e);
			if(cb) cb(e, null);
		}
	});
}

//rest call to peer to register
function rest_register(id, secret, cb){
	//console.log("register");
	var host = bag.peers[0].api_host;
	var port = bag.peers[0].api_port;
	var data = {
					enrollId: id,
					enrollSecret: secret
				};

	$.ajax({
		method: 'POST',
		url: 'https://' + host + ':' + port + '/registrar',
		data: JSON.stringify(data),
		contentType: 'application/json',
		success: function(json){
			console.log('Success - register', json);
			bag.peers[0].user = id;
			if(cb) cb(null, json);
		},
		error: function(e){
			console.log('Error - register', e);
			if(cb) cb(e, null);
		}
	});
}

//rest call to peer to get chain stats
function rest_chainstats(cb){
	//console.log("chainstats");
	var host = bag.peers[0].api_host;
	var port = bag.peers[0].api_port;

	$.ajax({
		method: 'GET',
		url: 'https://' + host + ':' + port + '/chain',
		contentType: 'application/json',
		success: function(json){
			json.height--;
			console.log('Success - chainstats', json);
			if(cb) cb(null, json);
		},
		error: function(e){
			console.log('Error - chainstats', e);
			if(cb) cb(e, null);
		}
	});
}

//rest call to peer to get stats for a block
function rest_blockstats(height, cb){
	//console.log("blockstats");
	var host = bag.peers[0].api_host;
	var port = bag.peers[0].api_port;

	$.ajax({
		method: 'GET',
		url: 'https://' + host + ':' + port + '/chain/blocks/' + height,
		contentType: 'application/json',
		success: function(json){
			console.log('Success - blockstats', json);
			json.height = height;
			new_block(json);
			if(cb) cb(null, json);
		},
		error: function(e){
			console.log('Error - blockstats', e);
			if(cb) cb(e, null);
		}
	});
}


// =================================================================================
// Other fun
// =================================================================================
//display seconds
function formatTime(ms){
	var elasped = Math.floor((Date.now() - ms*1000) / 1000);
	var str = '';
	var levels = 0;

	if(elasped >= 60*60*24){
		levels++;
		str =  Math.floor(elasped / (60*60*24)) + 'days ';
		elasped = elasped % (60*60*24);
	}
	if(elasped >= 60*60){
		levels++;
		if(levels < 2){
			str =  Math.floor(elasped / (60*60)) + 'hr ';
			elasped = elasped % (60*60);
		}
	}
	if(elasped >= 60){
		if(levels < 2){
			levels++;
			str +=  Math.floor(elasped / 60) + 'min ';
			elasped = elasped % 60;
		}
	}
	if(levels < 2){
		str +=  elasped + 'sec ';
	}

	return str;
}

//fancy date
function formatDate(date, fmt) {
	date = new Date(date);
	function pad(value) {
		return (value.toString().length < 2) ? '0' + value : value;
	}
	return fmt.replace(/%([a-zA-Z])/g, function (_, fmtCode) {
		var tmp;
		switch (fmtCode) {
		case 'Y':								//Year
			return date.getUTCFullYear();
		case 'M':								//Month 0 padded
			return pad(date.getUTCMonth() + 1);
		case 'd':								//Date 0 padded
			return pad(date.getUTCDate());
		case 'H':								//24 Hour 0 padded
			return pad(date.getUTCHours());
		case 'I':								//12 Hour 0 padded
			tmp = date.getUTCHours();
			if(tmp == 0) tmp = 12;				//00:00 should be seen as 12:00am
			else if(tmp > 12) tmp -= 12;
			return pad(tmp);
		case 'p':								//am / pm
			tmp = date.getUTCHours();
			if(tmp >= 12) return 'pm';
			return 'am';
		case 'P':								//AM / PM
			tmp = date.getUTCHours();
			if(tmp >= 12) return 'PM';
			return 'AM';
		case 'm':								//Minutes 0 padded
			return pad(date.getUTCMinutes());
		case 's':								//Seconds 0 padded
			return pad(date.getUTCSeconds());
		case 'r':								//Milliseconds 0 padded
			return pad(date.getUTCMilliseconds(), 3);
		case 'q':								//UTC timestamp
			return date.getTime();
		default:
			throw new Error('Unsupported format code: ' + fmtCode);
		}
	});
}
