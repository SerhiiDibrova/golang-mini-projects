package initializers

import grails.config.Config
import java.net.URL
import java.net.MalformedURLException

class ServiceUrlGenerator {
    private Config grailsApplicationConfig

    ServiceUrlGenerator(Config grailsApplicationConfig) {
        this.grailsApplicationConfig = grailsApplicationConfig
    }

    URL generateServiceUrl(String lat, String lon) {
        String serviceAddress = grailsApplicationConfig.depthService.url
        if (serviceAddress == null || serviceAddress.isEmpty()) {
            return null
        }

        String params = "lon:${lon};lat:${lat}"
        String serviceUrl = "${serviceAddress}?${params}"

        try {
            return new URL(serviceUrl)
        } catch (MalformedURLException e) {
            throw new RuntimeException("Error parsing service URL", e)
        }
    }
}