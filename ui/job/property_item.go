package job

import (
	"fmt"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type PropertyItem struct {
	Indent int

	Key    string
	Anchor string

	// HasDefaults shows if either this item or sub-items have defaults
	MissingValues bool

	Property *Property
}

func NewPropertyItems(props []Property) []PropertyItem {
	var items []*PropertyItem
	var itemsByIndent []*PropertyItem

	var prevProp *Property

	for i, prop := range props {
		depth, parts := matchingPropsDepth(prevProp, prop)

		lastJ := len(parts) - 1

		for j, part := range parts {
			if j < depth {
				continue
			} else {
				itemsByIndent = itemsByIndent[0:j]
			}

			item := PropertyItem{
				Indent: j,
				Key:    part,
			}

			if j == lastJ {
				item.Property = &props[i]
				item.MissingValues = !props[i].HasDefault()
				item.Anchor = prop.Name

				// Propagate missing value mark to the top level item
				if item.MissingValues {
					for _, item := range itemsByIndent {
						item.MissingValues = true
					}
				}
			} else {
				item.Anchor = strings.Join(parts[0:j+1], ".")
			}

			item.Anchor = "p=" + item.Anchor

			items = append(items, &item)
			itemsByIndent = append(itemsByIndent, &item)
		}

		prevProp = &props[i]
	}

	// Non-pointer items array
	copiedItems := []PropertyItem{}

	for _, item := range items {
		copiedItems = append(copiedItems, *item)
	}

	return copiedItems
}

func matchingPropsDepth(prevProp *Property, currProp Property) (int, []string) {
	currParts := strings.Split(currProp.Name, ".")

	if prevProp == nil {
		return 0, currParts
	}

	prevParts := strings.Split(prevProp.Name, ".")

	for i, currPart := range currParts {
		if len(prevParts) == i || prevParts[i] != currPart {
			return i, currParts
		}
	}

	return len(currParts), currParts
}

func (i PropertyItem) PoundAnchor() string {
	return "#" + i.Anchor
}

func (i PropertyItem) HasLongKey() bool {
	d, err := i.DefaultAsYAML()
	if err != nil {
		d = ""
	}

	parts := strings.Split(i.Key+": "+d, "\n")

	for _, v := range parts {
		if len(v) > 40 {
			return true
		}
	}

	return false
}

func (i PropertyItem) DefaultAsYAML() (string, error) {
	result, err := i.Property.DefaultAsYAML()
	if err != nil {
		return "", bosherr.WrapError(err, "Indenting property '%s' default", i.Property.Name)
	}

	// YAML encoder might add extra new line breaks
	result = strings.Trim(result, "\n")

	parts := strings.Split(result, "\n")

	if len(parts) == 1 {
		return result, nil
	}

	for j, v := range parts {
		// it does mean arrays and hashes will be indented
		parts[j] = fmt.Sprintf("  %s", v)
	}

	return "\n" + strings.Join(parts, "\n"), nil
}
