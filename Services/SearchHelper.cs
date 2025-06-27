using System;
using System.Collections.Generic;
using System.Windows.Controls;
using System.Windows.Documents;
using System.Windows.Media;

namespace CxExamClient.Services
{
    public class SearchHelper
    {
        private readonly List<TextRange> _searchResults = new();
        private int _currentSearchIndex = 0;

        public void HighlightText(RichTextBox richTextBox, string searchTerm)
        {
            if (string.IsNullOrEmpty(searchTerm) || richTextBox.Document == null)
                return;

            // 清除之前的高亮
            ClearHighlights(richTextBox);
            _searchResults.Clear();

            var document = richTextBox.Document;
            var textRange = new TextRange(document.ContentStart, document.ContentEnd);
            var text = textRange.Text;
            var startIndex = 0;

            while (true)
            {
                var index = text.IndexOf(searchTerm, startIndex, StringComparison.OrdinalIgnoreCase);
                if (index == -1) break;

                var start = GetTextPointerAtOffset(document.ContentStart, index);
                var end = GetTextPointerAtOffset(start, searchTerm.Length);
                
                if (start != null && end != null)
                {
                    var range = new TextRange(start, end);
                    _searchResults.Add(range);
                }

                startIndex = index + 1;
            }

            if (_searchResults.Count > 0)
            {
                _currentSearchIndex = 0;
                HighlightCurrentResult();
                
                // 滚动到第一个结果
                var firstResult = _searchResults[0];
                if (firstResult.Start != null)
                {
                    firstResult.Start.Paragraph?.BringIntoView();
                }
            }
        }

        public void NextResult(RichTextBox richTextBox)
        {
            if (_searchResults.Count == 0) return;

            _currentSearchIndex = (_currentSearchIndex + 1) % _searchResults.Count;
            HighlightCurrentResult();
            
            // 滚动到当前结果
            var currentResult = _searchResults[_currentSearchIndex];
            if (currentResult.Start != null)
            {
                currentResult.Start.Paragraph?.BringIntoView();
            }
        }

        private void HighlightCurrentResult()
        {
            // 清除所有高亮
            foreach (var result in _searchResults)
            {
                result.ApplyPropertyValue(TextElement.BackgroundProperty, Brushes.Transparent);
            }

            // 高亮当前结果
            if (_currentSearchIndex >= 0 && _currentSearchIndex < _searchResults.Count)
            {
                _searchResults[_currentSearchIndex].ApplyPropertyValue(TextElement.BackgroundProperty, Brushes.Yellow);
            }
        }

        private void ClearHighlights(RichTextBox richTextBox)
        {
            foreach (var result in _searchResults)
            {
                result.ApplyPropertyValue(TextElement.BackgroundProperty, Brushes.Transparent);
            }
        }

        private TextPointer? GetTextPointerAtOffset(TextPointer start, int offset)
        {
            var current = start;
            var currentOffset = 0;

            while (current != null && currentOffset < offset)
            {
                if (current.GetPointerContext(LogicalDirection.Forward) == TextPointerContext.Text)
                {
                    var textLength = current.GetTextRunLength(LogicalDirection.Forward);
                    if (currentOffset + textLength > offset)
                    {
                        return current.GetPositionAtOffset(offset - currentOffset);
                    }
                    currentOffset += textLength;
                }
                current = current.GetNextContextPosition(LogicalDirection.Forward);
            }

            return current;
        }
    }
} 