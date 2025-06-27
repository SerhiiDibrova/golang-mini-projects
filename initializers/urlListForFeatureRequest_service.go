package urlListForFeatureRequest

import grails.config.Config
import grails.core.support.GrailsApplication
import grails.web.servlet.mvc.GrailsHttpServletResponse
import grails.web.servlet.mvc.GrailsHttpServletRequest
import org.springframework.http.HttpStatus

class UrlListForFeatureRequestService {
    GrailsApplication grailsApplication

    def processRequest(GrailsHttpServletRequest request, GrailsHttpServletResponse response) {
        def urlListConfig = grailsApplication.config.marvl.urlList
        if (!urlListConfig) {
            response.setStatus(HttpStatus.INTERNAL_SERVER_ERROR.value())
            response.getWriter().write("URL list configuration not found")
            return
        }

        def fieldParam = requestSingleFieldParamProcessor.process(request)
        if (!fieldParam) {
            response.setStatus(HttpStatus.INTERNAL_SERVER_ERROR.value())
            response.getWriter().write("Failed to process field param")
            return
        }

        def stream = urlListStreamProcessor.process(urlListConfig, fieldParam)
        if (!stream) {
            response.setStatus(HttpStatus.INTERNAL_SERVER_ERROR.value())
            response.getWriter().write("Failed to process stream")
            return
        }

        def proxyingAllowed = _performProxyingIfAllowed(grailsApplication)
        if (proxyingAllowed) {
            stream = _performProxying(grailsApplication, stream)
            if (!stream) {
                response.setStatus(HttpStatus.INTERNAL_SERVER_ERROR.value())
                response.getWriter().write("Failed to perform proxying")
                return
            }
        }

        response.setContentType("application/json")
        response.getWriter().write(stream as JSON)
    }
}

class RequestSingleFieldParamProcessor {
    def process(GrailsHttpServletRequest request) {
        def fieldParam = request.getParameter("field")
        if (!fieldParam) {
            return null
        }
        fieldParam
    }
}

class UrlListStreamProcessor {
    def process(urlList, fieldParam) {
        def stream = []
        urlList.each { url ->
            stream << "${url}/${fieldParam}"
        }
        stream
    }
}

def _performProxyingIfAllowed(GrailsApplication grailsApplication) {
    def proxyingConfig = grailsApplication.config.marvl.proxyingEnabled
    if (!proxyingConfig) {
        return false
    }
    proxyingConfig
}

def _performProxying(GrailsApplication grailsApplication, stream) {
    // implement proxying logic here
    stream
}