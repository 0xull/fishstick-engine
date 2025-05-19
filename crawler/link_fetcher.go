package crawler

import (
	"context"

	"github.com/IkehAkinyemi/fishstick-engine/pipeline"
)

var _ pipeline.Processor = (*linkFetcher)(nil)

type linkFetcher struct {
	urlGetter URLGetter // TODO: Yemi, implement the URLGetter first
	netDetector PrivateNetworkDetector // TODO: likewise PrivateNetworkDetector!
}

func newLinkFetcher(urlGetter URLGetter, netDetector PrivateNetworkDetector) *linkFetcher {
	return &linkFetcher{
		urlGetter: urlGetter, 
		netDetector: netDetector,
	}
}

func (lf *linkFetcher) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {
	payload := p.(*crawlerPayload)
	
	// Skip URLs that point to files that cannot contain html content.
	// if exclusiveRegex.MatchString(payload.URL) // TODO: same for exclusiveRegex
	
	return payload, nil
}
