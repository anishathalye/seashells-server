(function() {
  'use strict';

  var term = new Terminal({
    convertEol: true,
    scrollback: 10000,
    disableStdin: true,
    cursorBlink: true,
  });
  term.open(document.getElementById('terminal'), {
    focus: true,
  });
  term.fit();
  window.addEventListener('resize', function() {
    term.fit();
  });

  var wsscheme;
  if (window.location.protocol.startsWith('https')) {
    wsscheme = 'wss://';
  } else {
    wsscheme = 'ws://';
  }
  var url = wsscheme + window.location.host + '/ws';

  var socket = new WebSocket(url);
  socket.binaryType = 'arraybuffer';

  socket.onmessage = function(event) {
    var data = event.data;
    var str = new TextDecoder('utf-8').decode(new Uint8Array(data));
    term.write(str);
  };

  socket.onopen = function() {
    socket.send(sessionId);
  };

  socket.onclose = function() {
    term.setOption('cursorBlink', false);
  };
}());
