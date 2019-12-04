package truncate_test

import (
	"fmt"

	"github.com/alextanhongpin/pkg/truncate"
)

func Example() {
	text := `The Creaking Pagoda or Chinese Summer-House is located in Tsarskoye Selo, outside Saint Petersburg, Russia, between two ponds on the boundary separating the Catherine Park of the Baroque Catherine Palace and the New Garden of the neoclassical Alexander Palace's Alexander Park. The pagoda, designed by Georg von Veldten, is a folly that resulted from the 18th-century taste for chinoiserie. The walls are decorated with figures of dragons and other stylized Chinese motifs. Construction lasted from 1778 to 1786, and the structure was restored from 1954 to 1956. The name of the structure refers to a characteristic sound produced by a metal weathervane, shaped like a banner, on the top of the structure, which creaks when it is turned by the wind.`

	fmt.Println(truncate.Paragraph(text, 200))
	fmt.Println(truncate.Sentence(text))
}
