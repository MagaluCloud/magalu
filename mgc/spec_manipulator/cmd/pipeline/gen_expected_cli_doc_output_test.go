package pipeline

import (
	"testing"
)

func Test_stripHtmlTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "texto sem tags HTML",
			input:    "Este é um texto simples sem tags",
			expected: "Este é um texto simples sem tags",
		},
		{
			name:     "texto com tag única",
			input:    "Este é um <b>texto</b> com negrito",
			expected: "Este é um texto com negrito",
		},
		{
			name:     "texto com múltiplas tags",
			input:    "Este é um <b>texto</b> com <i>itálico</i> e <u>sublinhado</u>",
			expected: "Este é um texto com itálico e sublinhado",
		},
		{
			name:     "string vazia",
			input:    "",
			expected: "",
		},
		{
			name:     "apenas tags HTML",
			input:    "<html><body><div></div></body></html>",
			expected: "",
		},
		{
			name:     "tags aninhadas",
			input:    "<div><p>Texto <span>aninhado</span> aqui</p></div>",
			expected: "Texto aninhado aqui",
		},
		{
			name:     "tags com atributos",
			input:    "<a href=\"https://example.com\" class=\"link\">Link</a>",
			expected: "Link",
		},
		{
			name:     "tags auto-fechadas",
			input:    "Linha 1<br/>Linha 2<hr/>Fim",
			expected: "Linha 1Linha 2Fim",
		},
		{
			name:     "texto com caracteres especiais HTML",
			input:    "Texto com &lt; e &gt; e &amp;",
			expected: "Texto com &lt; e &gt; e &amp;",
		},
		{
			name:     "tags mal formadas",
			input:    "Texto com <tag> não fechada e </tag> fechada",
			expected: "Texto com  não fechada e  fechada",
		},
		{
			name:     "quebras de linha e espaços",
			input:    "<p>\n  Texto com\n  quebras de linha\n</p>",
			expected: "\n  Texto com\n  quebras de linha\n",
		},
		{
			name:     "tags com números",
			input:    "<h1>Título 1</h1><h2>Título 2</h2>",
			expected: "Título 1Título 2",
		},
		{
			name:     "comentários HTML",
			input:    "Texto <!-- comentário --> aqui",
			expected: "Texto  aqui",
		},
		{
			name:     "tag com símbolos especiais",
			input:    "<tag-name data-value=\"test\">Conteúdo</tag-name>",
			expected: "Conteúdo",
		},
		{
			name:     "texto longo com várias tags",
			input:    "<html><head><title>Título</title></head><body><div class=\"container\"><p>Parágrafo com <strong>texto forte</strong> e <em>ênfase</em>.</p></div></body></html>",
			expected: "TítuloParágrafo com texto forte e ênfase.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripHtmlTags(tt.input)
			if result != tt.expected {
				t.Errorf("stripHtmlTags() = %q, esperado %q", result, tt.expected)
			}
		})
	}
}

func Benchmark_stripHtmlTags(t *testing.B) {
	input := "<html><head><title>Título de Teste</title></head><body><div class=\"container\"><p>Este é um <strong>parágrafo de teste</strong> com várias <em>tags HTML</em> para <a href=\"#\">benchmark</a>.</p><ul><li>Item 1</li><li>Item 2</li><li>Item 3</li></ul></div></body></html>"

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		stripHtmlTags(input)
	}
}

func Test_stripHtmlTags_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "tag incompleta no início",
			input:    "<incompleta Texto normal",
			expected: "<incompleta Texto normal",
		},
		{
			name:     "tag incompleta no final",
			input:    "Texto normal <incompleta",
			expected: "Texto normal <incompleta",
		},
		{
			name:     "múltiplos sinais de menor e maior",
			input:    "a < b > c < d",
			expected: "a  c < d",
		},
		{
			name:     "tags vazias",
			input:    "Texto<>com<>tags<>vazias",
			expected: "Textocomtagsvazias",
		},
		{
			name:     "apenas espaços dentro das tags",
			input:    "< >Texto< >aqui< >",
			expected: "Textoaqui",
		},
		{
			name:     "símbolos menor e maior isolados",
			input:    "a < b e c > d",
			expected: "a  d",
		},
		{
			name:     "tag válida entre símbolos isolados",
			input:    "a < b <span>texto</span> c > d",
			expected: "a texto c > d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripHtmlTags(tt.input)
			if result != tt.expected {
				t.Errorf("stripHtmlTags() = %q, esperado %q", result, tt.expected)
			}
		})
	}
}
