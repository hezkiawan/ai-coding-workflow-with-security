package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"opsdesk/internal/database"
	"opsdesk/internal/middleware"
	"opsdesk/internal/models"
	"opsdesk/internal/utils"
)

type XMLAssetItem struct {
	XMLName      xml.Name `xml:"Asset"`
	AssetTag     string   `xml:"AssetTag"`
	Name         string   `xml:"Name"`
	Category     string   `xml:"Category"`
	Model        string   `xml:"Model"`
	SerialNumber string   `xml:"SerialNumber"`
	Status       string   `xml:"Status"`
	Location     string   `xml:"Location"`
	IPAddress    string   `xml:"IPAddress"`
}

type XMLAssetBatch struct {
	XMLName xml.Name       `xml:"Assets"`
	Items   []XMLAssetItem `xml:"Asset"`
}

type JSONAssetBatch struct {
	Items []CreateAssetRequest `json:"items"`
}

func (h *AssetHandler) ImportXML(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetCurrentUser(r)
	if claims == nil {
		utils.Error(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Failed to read XML payload")
		return
	}

	var batch XMLAssetBatch
	if err := xml.Unmarshal(body, &batch); err != nil {
		utils.Error(w, http.StatusBadRequest, fmt.Sprintf("XML parsing error: %v", err))
		return
	}

	now := time.Now()
	importedCount := 0

	for _, item := range batch.Items {
		if item.AssetTag == "" || item.Name == "" {
			continue
		}
		status := item.Status
		if status == "" {
			status = "active"
		}

		_, err := database.DB.Exec(`
			INSERT OR REPLACE INTO assets (asset_tag, name, category, model, serial_number, status, location, ip_address, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, item.AssetTag, item.Name, item.Category, item.Model, item.SerialNumber, status, item.Location, item.IPAddress, now, now)

		if err == nil {
			importedCount++
		}
	}

	utils.Success(w, http.StatusOK, map[string]interface{}{
		"imported_count": importedCount,
		"total_parsed":   len(batch.Items),
	}, "XML assets batch import completed")
}

func (h *AssetHandler) ImportJSON(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetCurrentUser(r)
	if claims == nil {
		utils.Error(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var batch JSONAssetBatch
	if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid JSON payload structure")
		return
	}

	now := time.Now()
	importedCount := 0

	for _, item := range batch.Items {
		if item.AssetTag == "" || item.Name == "" || item.Category == "" {
			continue
		}
		status := item.Status
		if status == "" {
			status = models.AssetStatusActive
		}

		_, err := database.DB.Exec(`
			INSERT OR REPLACE INTO assets (asset_tag, name, category, model, serial_number, status, location, assigned_to_id, ip_address, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, item.AssetTag, item.Name, item.Category, item.Model, item.SerialNumber, status, item.Location, item.AssignedToID, item.IPAddress, now, now)

		if err == nil {
			importedCount++
		}
	}

	utils.Success(w, http.StatusOK, map[string]interface{}{
		"imported_count": importedCount,
		"total_parsed":   len(batch.Items),
	}, "JSON assets batch import completed")
}

func (h *AssetHandler) ExportReport(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	sortBy := r.URL.Query().Get("sort_by")
	search := r.URL.Query().Get("search")

	if sortBy == "" {
		sortBy = "id"
	}

	baseQuery := "SELECT id, asset_tag, name, category, model, serial_number, status, location, ip_address, created_at, updated_at FROM assets"
	if search != "" {
		baseQuery += fmt.Sprintf(" WHERE name LIKE '%%%s%%' OR category LIKE '%%%s%%'", search, search)
	}
	baseQuery += fmt.Sprintf(" ORDER BY %s ASC", sortBy)

	rows, err := database.DB.Query(baseQuery)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to generate asset export report")
		return
	}
	defer rows.Close()

	items := make([]models.Asset, 0)
	for rows.Next() {
		var a models.Asset
		if err := rows.Scan(
			&a.ID, &a.AssetTag, &a.Name, &a.Category, &a.Model, &a.SerialNumber,
			&a.Status, &a.Location, &a.IPAddress, &a.CreatedAt, &a.UpdatedAt,
		); err == nil {
			items = append(items, a)
		}
	}

	if format == "xml" {
		w.Header().Set("Content-Type", "application/xml")
		xmlBatch := XMLAssetBatch{}
		for _, a := range items {
			xmlBatch.Items = append(xmlBatch.Items, XMLAssetItem{
				AssetTag:     a.AssetTag,
				Name:         a.Name,
				Category:     a.Category,
				Model:        a.Model,
				SerialNumber: a.SerialNumber,
				Status:       string(a.Status),
				Location:     a.Location,
				IPAddress:    a.IPAddress,
			})
		}
		_ = xml.NewEncoder(w).Encode(xmlBatch)
		return
	}

	utils.Success(w, http.StatusOK, items)
}
