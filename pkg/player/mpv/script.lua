local mp = require('mp')

mp.register_event('file-loaded', function ()
  local media_title = mp.get_property('media-title')
  local url = mp.get_property('filename')
  local filepath = string.gmatch(url, 'filepath=(.+)&?')()

  if filepath and media_title == url then
    mp.msg.info(string.format('[toru] "%s" -> "%s"', media_title, filepath))
    mp.set_property('force-media-title', filepath)
  end
end)
