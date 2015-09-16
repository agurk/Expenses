#!/usr/bin/python

import config
import os.path
import Image

class ImageHandler:

    def Image(self, filename, path):
        if path == 'thumbnail':
            return self._Thumbnail(filename)
        else:
            if os.path.isfile(config.IMAGE_DIR + '/' + filename):
                return config.IMAGE_DIR + '/' + filename
            else:
                return config.IMAGE_DIR + '/' + config.IMAGE_404

    def _Thumbnail(self, filename):
        if not os.path.isfile(config.IMAGE_DIR + '/' + filename):
            return self._Thumbnail(config.IMAGE_404)
        success = 1
        if not os.path.isfile(config.THUMBNAIL_CACHE + '/' + filename):
            success = self._CreateThumbnail(filename)
        if success:
            return config.THUMBNAIL_CACHE + '/' + filename
        else:
            return config.IMAGE_404
        
    def _CreateThumbnail(self, filename):
        try:
            im = Image.open(config.IMAGE_DIR + '/' + filename)
            im.thumbnail(config.THUMBNAIL_SIZE, Image.ANTIALIAS)
            im.save(config.THUMBNAIL_CACHE + '/' + filename, config.THUMBNAIL_FORMAT)
        except:
            print 'Error converting: ' + filename
            return 0
        return 1

