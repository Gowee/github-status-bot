/// A wrapper around the page API of statuspage

use anyhow::Result;

struct StatusPage {
    page_id: String,
}

impl StatusPage {
    fn new(page_id: impl AsRef<str>) -> StatusPage {
        StatusPage {
            page_id: page_id.as_ref().to_owned(),
        }
    }


    async fn get_summary(&self) -> Result<Summary> {
        reqwest::get(format!(
            "https://{}.statuspage.io/api/v2/summary.json",
            self.page_id
        ))
        .await?
        .json::<Summary>()
        .await?
    }
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
pub struct Summary {
    pub page: Page,
    pub components: Vec<Component>,
    pub incidents: Vec<::serde_json::Value>,
    pub scheduled_maintenances: Vec<::serde_json::Value>,
    pub status: Status,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
pub struct Page {
    pub id: String,
    pub name: String,
    pub url: String,
    pub time_zone: String,
    pub updated_at: String,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
pub struct Component {
    pub id: String,
    pub name: String,
    pub status: String,
    pub created_at: String,
    pub updated_at: String,
    pub position: i64,
    pub description: Option<String>,
    pub showcase: bool,
    pub start_date: ::serde_json::Value,
    pub group_id: ::serde_json::Value,
    pub page_id: String,
    pub group: bool,
    pub only_show_if_degraded: bool,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
pub struct Status {
    pub indicator: String,
    pub description: String,
}
