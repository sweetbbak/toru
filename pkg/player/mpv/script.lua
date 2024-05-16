local mp = require('mp')

mp.add_hook('on_load', 50, function()
  local media_title = mp.get_property('media-title')
  local filename = mp.get_property('filename')
  local external_title = mp.get_opt('external-title')

  if media_title == filename then
    mp.msg.info(media_title)
    mp.msg.info(filename)
    mp.msg.info(external_title)
    mp.set_property('force-media-title', external_title)
  end
end)
