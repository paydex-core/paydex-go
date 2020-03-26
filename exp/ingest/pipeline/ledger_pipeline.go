package pipeline

import (
	"github.com/paydex-core/paydex-go/exp/ingest/io"
	supportPipeline "github.com/paydex-core/paydex-go/exp/support/pipeline"
)

func LedgerNode(processor LedgerProcessor) *supportPipeline.PipelineNode {
	return &supportPipeline.PipelineNode{
		Processor: &ledgerProcessorWrapper{processor},
	}
}

func (p *LedgerPipeline) Process(reader io.LedgerReader) <-chan error {
	return p.Pipeline.Process(&ledgerReaderWrapper{reader})
}
