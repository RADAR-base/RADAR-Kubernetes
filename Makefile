.PHONY: gen
gen: gen-chart-doc

.PHONY: gen-chart-doc
gen-chart-doc:
	@echo "Generate chart docs"
	@helm-docs -s file --template-files=charts/_templates.gotmpl --template-files=DOCS.md.gotmpl --template-files=README.md.gotmpl -c charts/management-portal
