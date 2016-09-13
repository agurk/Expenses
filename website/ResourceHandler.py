import config
import os.path
from PIL import Image

from enum import Enum
class ResourceType(Enum):
    image = 1
    pdf = 2

class ResourceHandler:

    def Resource(self, resType, filename, path):
        print (resType)
        if path == 'thumbnail':
            return self._Thumbnail(resType, filename)
        elif resType == ResourceType.image:
            return self.Image(filename)
        elif resType == ResourceType.pdf:
            return self.PDF(filename)

    def PDF(self, filename):
       if os.path.isfile(config.PDF_DIR + '/' + filename):
           return config.PDF_DIR + '/' + filename
       else:
           return config.IMAGE_DIR + '/' + config.IMAGE_404

    def _PDFThumbnail(self, filename):
        if not os.path.isfile(config.PDF_DIR + '/' + filename):
            return self._Thumbnail(config.IMAGE_404)
        success = 1
        if not os.path.isfile(config.THUMBNAIL_CACHE + '/' + filename):
            success = self._CreateThumbnail(filename)
        if success:
            return config.THUMBNAIL_CACHE + '/' + filename
        else:
            return config.IMAGE_404

    def Image(self, filename):
       if os.path.isfile(config.IMAGE_DIR + '/' + filename):
           return config.IMAGE_DIR + '/' + filename
       else:
           return config.IMAGE_DIR + '/' + config.IMAGE_404

    def _Thumbnail(self, filetype, filename):
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
            print ('Error converting: ' + filename)
            return 0
        return 1

